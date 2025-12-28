#!/usr/bin/env python3
import json
import yaml
import re
import os
import sys
from typing import Dict, List, Any, Tuple, Optional
import shutil

def load_yaml_file(filepath: str) -> Optional[Dict]:
    """加载YAML文件"""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            return yaml.safe_load(f)
    except Exception as e:
        print(f"✗ 加载YAML文件失败: {e}")
        return None

def save_yaml_file(data: Dict, filepath: str) -> bool:
    """保存YAML文件"""
    try:
        with open(filepath, 'w', encoding='utf-8') as f:
            yaml.dump(data, f, allow_unicode=True, sort_keys=False, default_flow_style=False)
        return True
    except Exception as e:
        print(f"✗ 保存YAML文件失败: {e}")
        return False

def load_json_file(filepath: str) -> Optional[Dict]:
    """加载JSON文件"""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            return json.load(f)
    except Exception as e:
        print(f"✗ 加载JSON文件失败: {e}")
        return None

def save_json_file(data: Dict, filepath: str) -> bool:
    """保存JSON文件"""
    try:
        with open(filepath, 'w', encoding='utf-8') as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        return True
    except Exception as e:
        print(f"✗ 保存JSON文件失败: {e}")
        return False

def fix_swagger_object(data: Dict) -> Tuple[Dict, Dict]:
    """
    修复swagger对象中的example字段

    Args:
        data: swagger对象

    Returns:
        Tuple[修复后的对象, 统计信息]
    """
    if not isinstance(data, dict):
        return data, {}

    stats = {
        "total_params": 0,
        "params_fixed": 0,
        "other_examples": 0,
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

                # 记录参数类型
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
    fixed_data = fix_param_examples(data.copy())
    return fixed_data, stats

def process_swagger_file(filepath: str, output_path: str = None, dry_run: bool = False) -> bool:
    """
    处理单个swagger文件

    Args:
        filepath: 输入文件路径
        output_path: 输出文件路径，None则使用输入路径
        dry_run: 是否只显示统计不实际修改

    Returns:
        bool: 是否成功
    """
    if not os.path.exists(filepath):
        print(f"✗ 文件不存在: {filepath}")
        return False

    filename = os.path.basename(filepath)
    file_ext = os.path.splitext(filename)[1].lower()

    print(f"\n处理文件: {filename}")
    print("-" * 50)

    # 加载数据
    data = None
    is_go_file = filename.lower().endswith('.go')

    if is_go_file:
        # 处理Go文件
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()

        # 查找反引号内的JSON内容
        go_matches = re.findall(r'`({[\s\S]*?})`', content)
        if not go_matches:
            print("✗ 在Go文件中未找到JSON内容")
            return False

        json_str = go_matches[0]
        try:
            data = json.loads(json_str)
        except json.JSONDecodeError as e:
            print(f"✗ 解析Go文件中的JSON失败: {e}")
            return False

    elif file_ext in ['.yaml', '.yml']:
        data = load_yaml_file(filepath)
    elif file_ext == '.json':
        data = load_json_file(filepath)
    else:
        print(f"✗ 不支持的文件类型: {file_ext}")
        return False

    if data is None:
        return False

    # 修复数据
    fixed_data, stats = fix_swagger_object(data)

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

    success = False
    if is_go_file:
        # 保存Go文件
        fixed_json = json.dumps(fixed_data, ensure_ascii=False, indent=2)
        # 替换原JSON内容
        new_content = re.sub(r'`({[\s\S]*?})`', f'`{fixed_json}`', content, count=1)

        try:
            with open(output_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
            print(f"✓ 已保存到: {output_path}")
            success = True
        except Exception as e:
            print(f"✗ 保存Go文件失败: {e}")
            success = False

    elif file_ext in ['.yaml', '.yml']:
        success = save_yaml_file(fixed_data, output_path)
        if success:
            print(f"✓ 已保存到: {output_path}")

    elif file_ext == '.json':
        success = save_json_file(fixed_data, output_path)
        if success:
            print(f"✓ 已保存到: {output_path}")

    return success

def process_multiple_files(file_paths: List[str], output_dir: str = None, dry_run: bool = False) -> Dict:
    """
    批量处理多个文件

    Args:
        file_paths: 文件路径列表
        output_dir: 输出目录，None则覆盖原文件
        dry_run: 是否只显示统计不实际修改

    Returns:
        处理结果统计
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

        success = process_swagger_file(filepath, output_path, dry_run)

        if success:
            results["success"] += 1
            results["details"].append({"file": filepath, "status": "success"})
        else:
            results["failed"] += 1
            results["details"].append({"file": filepath, "status": "failed"})

    return results

def find_swagger_files(directory: str, patterns: List[str] = None) -> List[str]:
    """
    查找目录下的swagger相关文件

    Args:
        directory: 目录路径
        patterns: 文件模式列表

    Returns:
        找到的文件路径列表
    """
    if patterns is None:
        patterns = [
            "**/*.json",
            "**/*.yaml",
            "**/*.yml",
            "**/docs.go",
            "**/swagger.go"
        ]

    import glob

    found_files = []
    for pattern in patterns:
        files = glob.glob(os.path.join(directory, pattern), recursive=True)
        found_files.extend(files)

    # 去重
    found_files = list(set(found_files))

    # 过滤出可能是swagger文档的文件
    swagger_files = []
    for filepath in found_files:
        filename = os.path.basename(filepath).lower()

        # 检查文件内容是否包含swagger相关关键词
        try:
            with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read(4096)  # 只读取前4KB

            # 检查是否包含swagger相关关键词
            swagger_keywords = ['swagger', 'openapi', '"paths"', '"parameters"']
            if any(keyword.lower() in content.lower() for keyword in swagger_keywords):
                swagger_files.append(filepath)
        except:
            continue

    return swagger_files

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
        print("Swagger文档修复工具")
        print("=" * 60)
        print("修复swagger文档中parameters的example字段为x-example")
        print()
        print("用法:")
        print("  1. 处理单个文件:")
        print("     python fix_swagger_all.py <文件路径>")
        print()
        print("  2. 批量处理目录下所有相关文件:")
        print("     python fix_swagger_all.py <目录路径> --batch")
        print()
        print("  3. 模拟运行（不实际修改）:")
        print("     python fix_swagger_all.py <文件或目录> --dry-run")
        print()
        print("  4. 指定输出目录:")
        print("     python fix_swagger_all.py <输入目录> <输出目录>")
        print()
        print("  5. 处理多个指定文件:")
        print("     python fix_swagger_all.py 文件1 文件2 文件3")
        print()
        print("示例:")
        print("  python fix_swagger_all.py swagger.json")
        print("  python fix_swagger_all.py swagger.yaml")
        print("  python fix_swagger_all.py docs.go")
        print("  python fix_swagger_all.py api/ --batch")
        print("  python fix_swagger_all.py swagger.json --dry-run")
        print("  python fix_swagger_all.py swagger.yaml docs.go")
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
        elif arg.startswith("--output="):
            output_dir = arg.split("=", 1)[1]
        elif arg == "--output" and i + 1 < len(sys.argv):
            output_dir = sys.argv[i + 1]
            i += 1
        elif os.path.isdir(arg):
            if batch_mode:
                # 批量处理目录
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
            # 可能是文件但不存在
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
        for filepath in file_paths:
            if os.path.exists(filepath):
                backup_file(filepath)

    # 处理文件
    print("\n开始处理...")
    results = process_multiple_files(file_paths, output_dir, dry_run)

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

    if dry_run:
        print("\n提示: 使用 --backup 参数可以在修改前备份原文件")
        print("      示例: python fix_swagger_all.py swagger.json --backup")

if __name__ == "__main__":
    main()