#!/usr/bin/env python3
import re
import os
import sys
import shutil
import glob
from typing import Dict, List, Any, Optional
import json
import yaml

def process_go_file_simple(filepath: str, output_path: str = None, dry_run: bool = False) -> bool:
    """
    简单处理Go文件中的Swagger文档
    只修改parameters数组中的example字段
    """
    print(f"\n处理Go文件: {os.path.basename(filepath)}")
    print("-" * 50)

    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        print(f"✗ 读取文件失败: {e}")
        return False

    # 查找所有的parameters数组
    # 匹配模式："parameters": [ ... ]
    # 使用非贪婪匹配
    param_pattern = r'(\"parameters\"\s*:\s*\[)(.*?)(\][\s]*,)'

    # 编译正则表达式，使用 DOTALL 标志以匹配多行
    param_regex = re.compile(param_pattern, re.DOTALL)

    matches = list(param_regex.finditer(content))

    if not matches:
        print("✗ 未找到parameters数组")
        return False

    print(f"✓ 找到 {len(matches)} 个parameters数组")

    # 统计信息
    stats = {
        "total_params": 0,
        "params_fixed": 0,
        "path_params": 0,
        "query_params": 0,
        "header_params": 0,
        "cookie_params": 0
    }

    new_content = content
    offset = 0  # 由于我们逐个替换，需要跟踪偏移量

    for match in matches:
        full_match_start = match.start(0) + offset
        full_match_end = match.end(0) + offset

        # 获取整个匹配的字符串
        full_match = match.group(0)

        # 获取数组内容
        array_content = match.group(2)

        # 在数组内容中查找参数对象
        # 参数对象通常以 { 开始，以 } 结束
        param_objects = []

        # 使用简单的方法：查找所有 { ... } 对
        # 假设没有嵌套的对象
        stack = []
        current_start = None

        for i, char in enumerate(array_content):
            if char == '{':
                if not stack:
                    current_start = i
                stack.append('{')
            elif char == '}':
                if stack:
                    stack.pop()
                    if not stack and current_start is not None:
                        param_objects.append((current_start, i + 1))
                        current_start = None

        print(f"  - 在数组中找到 {len(param_objects)} 个参数对象")

        fixed_array_content = array_content
        param_offset = 0

        for start, end in param_objects:
            param_str = array_content[start:end]

            # 检查是否是参数对象（包含 "in" 字段）
            if '"in"' in param_str:
                # 提取参数位置
                in_match = re.search(r'"in"\s*:\s*"(\w+)"', param_str)
                if in_match:
                    param_in = in_match.group(1)
                    stats["total_params"] += 1

                    if param_in == "path":
                        stats["path_params"] += 1
                    elif param_in == "query":
                        stats["query_params"] += 1
                    elif param_in == "header":
                        stats["header_params"] += 1
                    elif param_in == "cookie":
                        stats["cookie_params"] += 1

                    # 检查是否有example字段
                    if '"example"' in param_str:
                        # 替换example为x-example
                        new_param_str = param_str.replace('"example"', '"x-example"')

                        # 更新数组内容
                        actual_start = start + param_offset
                        actual_end = end + param_offset
                        fixed_array_content = (
                            fixed_array_content[:actual_start] +
                            new_param_str +
                            fixed_array_content[actual_end:]
                        )

                        # 更新偏移量（因为字符串长度可能变化，但这里长度相同）
                        param_offset += len(new_param_str) - len(param_str)
                        stats["params_fixed"] += 1

        # 重建完整的匹配
        new_full_match = match.group(1) + fixed_array_content + match.group(3)

        # 替换原内容
        new_content = (
            new_content[:full_match_start] +
            new_full_match +
            new_content[full_match_end:]
        )

        # 更新偏移量
        offset += len(new_full_match) - len(full_match)

    # 显示统计信息
    print(f"统计信息:")
    print(f"  总参数数量: {stats.get('total_params', 0)}")
    print(f"  已修复的参数: {stats.get('params_fixed', 0)}")

    if stats.get('path_params', 0) > 0:
        print(f"  - path参数: {stats.get('path_params', 0)}")
    if stats.get('query_params', 0) > 0:
        print(f"  - query参数: {stats.get('query_params', 0)}")
    if stats.get('header_params', 0) > 0:
        print(f"  - header参数: {stats.get('header_params', 0)}")
    if stats.get('cookie_params', 0) > 0:
        print(f"  - cookie参数: {stats.get('cookie_params', 0)}")

    if dry_run:
        print("✓ 模拟运行完成，未实际修改文件")
        return True

    # 保存文件
    output_path = output_path or filepath
    try:
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"✓ 已保存到: {output_path}")
        return True
    except Exception as e:
        print(f"✗ 保存文件失败: {e}")
        return False

def process_swagger_file_simple(filepath: str, output_path: str = None, dry_run: bool = False) -> bool:
    """
    简单处理单个swagger文件
    只修改parameters数组中的example字段
    """
    if not os.path.exists(filepath):
        print(f"✗ 文件不存在: {filepath}")
        return False

    filename = os.path.basename(filepath)
    file_ext = os.path.splitext(filename)[1].lower()

    # 判断是否是Go文件
    is_go_file = filename.lower().endswith('.go')

    if is_go_file:
        return process_go_file_simple(filepath, output_path, dry_run)

    # 对于非Go文件，可以使用原有的JSON/YAML解析方法
    print(f"\n处理文件: {filename}")
    print("-" * 50)

    # 加载数据
    data = None

    if file_ext in ['.yaml', '.yml']:
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                data = yaml.safe_load(f)
        except Exception as e:
            print(f"✗ 加载YAML文件失败: {e}")
            return False
    elif file_ext == '.json':
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                data = json.load(f)
        except Exception as e:
            print(f"✗ 加载JSON文件失败: {e}")
            return False
    else:
        print(f"✗ 不支持的文件类型: {file_ext}")
        return False

    if data is None:
        return False

    # 修复数据
    stats = {
        "total_params": 0,
        "params_fixed": 0,
        "path_params": 0,
        "query_params": 0,
        "header_params": 0,
        "cookie_params": 0
    }

    def fix_param_examples(obj: Any) -> Any:
        """修复参数中的example字段"""
        if isinstance(obj, dict):
            # 如果是参数且有example，则修复
            if "in" in obj and obj["in"] in ["path", "query", "header", "cookie"]:
                stats["total_params"] += 1

                param_in = obj.get("in", "unknown")
                if param_in == "path":
                    stats["path_params"] += 1
                elif param_in == "query":
                    stats["query_params"] += 1
                elif param_in == "header":
                    stats["header_params"] += 1
                elif param_in == "cookie":
                    stats["cookie_params"] += 1

                # 修复example字段
                if "example" in obj:
                    obj["x-example"] = obj.pop("example")
                    stats["params_fixed"] += 1

            # 递归处理
            for key, value in obj.items():
                if key != "x-example":  # 避免递归新添加的字段
                    obj[key] = fix_param_examples(value)

        elif isinstance(obj, list):
            for i, item in enumerate(obj):
                obj[i] = fix_param_examples(item)

        return obj

    # 修复整个对象
    fixed_data = fix_param_examples(data)

    # 显示统计信息
    print(f"统计信息:")
    print(f"  总参数数量: {stats.get('total_params', 0)}")
    print(f"  已修复的参数: {stats.get('params_fixed', 0)}")

    if stats.get('path_params', 0) > 0:
        print(f"  - path参数: {stats.get('path_params', 0)}")
    if stats.get('query_params', 0) > 0:
        print(f"  - query参数: {stats.get('query_params', 0)}")
    if stats.get('header_params', 0) > 0:
        print(f"  - header参数: {stats.get('header_params', 0)}")
    if stats.get('cookie_params', 0) > 0:
        print(f"  - cookie参数: {stats.get('cookie_params', 0)}")

    if dry_run:
        print("✓ 模拟运行完成，未实际修改文件")
        return True

    # 保存数据
    if output_path is None:
        output_path = filepath

    try:
        if file_ext in ['.yaml', '.yml']:
            with open(output_path, 'w', encoding='utf-8') as f:
                yaml.dump(fixed_data, f, allow_unicode=True, sort_keys=False, default_flow_style=False)
        elif file_ext == '.json':
            with open(output_path, 'w', encoding='utf-8') as f:
                json.dump(fixed_data, f, ensure_ascii=False, indent=2)

        print(f"✓ 已保存到: {output_path}")
        return True
    except Exception as e:
        print(f"✗ 保存文件失败: {e}")
        return False

def process_multiple_files_simple(file_paths: List[str], output_dir: str = None, dry_run: bool = False) -> Dict:
    """
    批量处理多个文件
    """
    results = {
        "total": len(file_paths),
        "success": 0,
        "failed": 0,
        "details": []
    }

    for filepath in file_paths:
        if not os.path.exists(filepath):
            print(f"\n✗ 文件不存在: {filepath}")
            results["failed"] += 1
            results["details"].append({"file": filepath, "status": "not_found"})
            continue

        # 确定输出路径
        output_path = None
        if output_dir:
            filename = os.path.basename(filepath)
            output_path = os.path.join(output_dir, filename)

        success = process_swagger_file_simple(filepath, output_path, dry_run)

        if success:
            results["success"] += 1
            results["details"].append({"file": filepath, "status": "success"})
        else:
            results["failed"] += 1
            results["details"].append({"file": filepath, "status": "failed"})

    return results

def find_swagger_files(directory: str) -> List[str]:
    """查找目录下的swagger相关文件"""
    patterns = [
        "**/*swagger*.json",
        "**/*openapi*.json",
        "**/*swagger*.yaml",
        "**/*openapi*.yaml",
        "**/*swagger*.yml",
        "**/*openapi*.yml",
        "**/docs.go",
        "**/swagger*.go"
    ]

    found_files = []
    for pattern in patterns:
        try:
            files = glob.glob(os.path.join(directory, pattern), recursive=True)
            found_files.extend(files)
        except:
            continue

    # 去重
    found_files = list(set(found_files))

    return found_files

def backup_file(filepath: str) -> bool:
    """备份文件"""
    if not os.path.exists(filepath):
        return False

    backup_path = f"{filepath}.bak"
    counter = 1
    while os.path.exists(backup_path):
        backup_path = f"{filepath}.bak.{counter}"
        counter += 1

    try:
        shutil.copy2(filepath, backup_path)
        print(f"  ✓ 已备份: {backup_path}")
        return True
    except Exception as e:
        print(f"  ✗ 备份失败: {e}")
        return False

def main():
    """主函数"""
    if len(sys.argv) < 2:
        print("Swagger文档修复工具 (简单版)")
        print("=" * 60)
        print("修复swagger文档中parameters的example字段为x-example")
        print("专门处理Go模板文件，使用字符串替换方法")
        print()
        print("用法:")
        print("  python fix_swagger_simple.py <文件路径> [--dry-run] [--backup]")
        print("  python fix_swagger_simple.py <目录> --batch")
        print()
        print("示例:")
        print("  python fix_swagger_simple.py docs.go")
        print("  python fix_swagger_simple.py docs.go --dry-run")
        print("  python fix_swagger_simple.py docs.go --backup")
        print("  python fix_swagger_simple.py . --batch")
        sys.exit(1)

    # 解析参数
    file_paths = []
    output_dir = None
    dry_run = False
    batch_mode = False
    backup_files = False

    i = 1
    while i < len(sys.argv):
        arg = sys.argv[i]

        if arg == "--dry-run":
            dry_run = True
        elif arg == "--batch":
            batch_mode = True
        elif arg == "--backup":
            backup_files = True
        elif arg in ["--output", "-o"] and i + 1 < len(sys.argv):
            output_dir = sys.argv[i + 1]
            i += 1
        elif os.path.isdir(arg):
            if batch_mode:
                found_files = find_swagger_files(arg)
                if found_files:
                    file_paths.extend(found_files)
                else:
                    print(f"✗ 在目录 {arg} 中未找到swagger相关文件")
            else:
                print(f"✗ 错误: 目录 {arg} 需要与 --batch 参数一起使用")
                sys.exit(1)
        elif os.path.isfile(arg):
            file_paths.append(arg)
        elif not arg.startswith("-"):
            print(f"✗ 警告: 文件或目录不存在: {arg}")

        i += 1

    # 如果没有指定文件，检查当前目录
    if not file_paths and batch_mode:
        found_files = find_swagger_files(".")
        if found_files:
            file_paths = found_files
        else:
            print("✗ 在当前目录中未找到swagger相关文件")
            sys.exit(1)
    elif not file_paths:
        print("✗ 错误: 未指定要处理的文件")
        sys.exit(1)

    # 显示要处理的文件
    print(f"找到 {len(file_paths)} 个文件:")
    for filepath in file_paths:
        print(f"  - {filepath}")

    if dry_run:
        print("\n" + "="*60)
        print("模拟运行模式 (不会实际修改文件)")
        print("="*60)

    # 备份文件
    if backup_files and not dry_run:
        print("\n备份文件中...")
        backup_count = 0
        for filepath in file_paths:
            if os.path.exists(filepath):
                if backup_file(filepath):
                    backup_count += 1
        print(f"✓ 已备份 {backup_count} 个文件")

    # 处理文件
    print("\n开始处理...")
    results = process_multiple_files_simple(file_paths, output_dir, dry_run)

    # 显示结果
    print("\n" + "="*60)
    print("处理完成!")
    print(f"  成功: {results['success']}")
    print(f"  失败: {results['failed']}")
    print(f"  总计: {results['total']}")

    if results['failed'] > 0:
        print("\n失败的文件:")
        for detail in results['details']:
            if detail['status'] == 'failed':
                print(f"  - {detail['file']}")

if __name__ == "__main__":
    main()