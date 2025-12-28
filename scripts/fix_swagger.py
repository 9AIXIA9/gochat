#!/usr/bin/env python3
import json
import sys
import os
import re
from typing import Dict, List, Any, Optional

def fix_swagger_selective(input_file: str, output_file: str = None, dry_run: bool = False) -> Dict:
    """
    选择性修复swagger文档中的example字段

    只处理以下位置的example，改为x-example：
    1. 所有parameters中的example（无论in是path、query、header还是cookie）
    2. 不处理其他地方（schema、responses等）的example

    Args:
        input_file: 输入的swagger文件路径
        output_file: 输出文件路径，默认在原文件名后加_fixed
        dry_run: 只统计不修改

    Returns:
        统计信息字典
    """
    if output_file is None and not dry_run:
        base, ext = os.path.splitext(input_file)
        output_file = f"{base}_fixed{ext}"

    try:
        with open(input_file, 'r', encoding='utf-8') as f:
            data = json.load(f)

        # 统计信息
        stats = {
            "total_params": 0,
            "params_fixed": 0,
            "other_examples": 0,
            "path_params": 0,
            "query_params": 0,
            "header_params": 0,
            "cookie_params": 0
        }

        def is_parameter_example(obj: Dict, path: str) -> bool:
            """
            判断一个example字段是否在parameter中
            参数定义通常包含以下字段：name, in, description, required, type等
            """
            # 如果对象是参数，它应该有"in"字段
            if isinstance(obj, dict) and "in" in obj:
                return obj["in"] in ["path", "query", "header", "cookie"]
            return False

        def process_parameters(params: List[Dict], context: str = "") -> List[Dict]:
            """处理parameters数组"""
            if not isinstance(params, list):
                return params

            fixed_params = []
            for param in params:
                if not isinstance(param, dict):
                    fixed_params.append(param)
                    continue

                param_copy = param.copy()
                stats["total_params"] += 1

                # 记录参数类型
                param_in = param_copy.get("in", "unknown")
                if param_in == "path":
                    stats["path_params"] += 1
                elif param_in == "query":
                    stats["query_params"] += 1
                elif param_in == "header":
                    stats["header_params"] += 1
                elif param_in == "cookie":
                    stats["cookie_params"] += 1

                # 如果parameter中有example，则修改
                if "example" in param_copy and is_parameter_example(param_copy, f"{context}.parameters"):
                    # 在修改前记录
                    param_name = param_copy.get("name", "unnamed")
                    param_example = param_copy["example"]

                    if dry_run:
                        print(f"  [DRY RUN] 将修改: {context} -> {param_name} ({param_in})")
                        print(f"    当前example: {param_example}")
                    else:
                        param_copy["x-example"] = param_copy.pop("example")

                    stats["params_fixed"] += 1

                fixed_params.append(param_copy)

            return fixed_params

        def process_paths(paths: Dict) -> Dict:
            """处理paths部分"""
            if not isinstance(paths, dict):
                return paths

            fixed_paths = {}
            for path, methods in paths.items():
                if not isinstance(methods, dict):
                    fixed_paths[path] = methods
                    continue

                fixed_methods = {}
                for method, operation in methods.items():
                    if not isinstance(operation, dict):
                        fixed_methods[method] = operation
                        continue

                    operation_copy = operation.copy()
                    context = f"paths.{path}.{method}"

                    # 处理operation的parameters
                    if "parameters" in operation:
                        operation_copy["parameters"] = process_parameters(
                            operation["parameters"],
                            context
                        )

                    # 处理请求体
                    if "requestBody" in operation and "content" in operation["requestBody"]:
                        # 不修改请求体中的example
                        pass

                    # 处理响应
                    if "responses" in operation:
                        # 不修改响应中的example
                        pass

                    fixed_methods[method] = operation_copy

                fixed_paths[path] = fixed_methods

            return fixed_paths

        def process_components(components: Dict) -> Dict:
            """处理components部分"""
            if not isinstance(components, dict):
                return components

            components_copy = components.copy()

            # 处理components/parameters
            if "parameters" in components and isinstance(components["parameters"], dict):
                params_dict = components["parameters"]
                fixed_params = {}

                for param_name, param_def in params_dict.items():
                    if isinstance(param_def, dict):
                        param_copy = param_def.copy()
                        if "example" in param_copy and is_parameter_example(param_copy, f"components.parameters.{param_name}"):
                            param_in = param_copy.get("in", "unknown")
                            if dry_run:
                                print(f"  [DRY RUN] 将修改: components.parameters.{param_name} ({param_in})")
                                print(f"    当前example: {param_copy['example']}")
                            else:
                                param_copy["x-example"] = param_copy.pop("example")
                            stats["params_fixed"] += 1
                        fixed_params[param_name] = param_copy
                    else:
                        fixed_params[param_name] = param_def

                components_copy["parameters"] = fixed_params

            return components_copy

        def find_other_examples(obj: Any, path: str = "") -> None:
            """查找但不修改其他位置的example，用于统计"""
            if isinstance(obj, dict):
                for key, value in obj.items():
                    if key == "example":
                        # 检查父对象是否是parameter
                        parent_is_param = False
                        if "in" in obj and obj["in"] in ["path", "query", "header", "cookie"]:
                            parent_is_param = True

                        if not parent_is_param:
                            stats["other_examples"] += 1
                            if dry_run and stats["other_examples"] <= 5:  # 只显示前5个
                                example_preview = str(value)[:50] + "..." if len(str(value)) > 50 else str(value)
                                print(f"  [DRY RUN] 保持原样: {path}.example ({example_preview})")

                    # 继续递归查找
                    if isinstance(value, (dict, list)):
                        new_path = f"{path}.{key}" if path else key
                        find_other_examples(value, new_path)

            elif isinstance(obj, list):
                for i, item in enumerate(obj):
                    if isinstance(item, (dict, list)):
                        new_path = f"{path}[{i}]" if path else f"[{i}]"
                        find_other_examples(item, new_path)

        # 主要处理逻辑
        print(f"开始处理文件: {input_file}")

        # 1. 处理根级别parameters
        if "parameters" in data and isinstance(data["parameters"], list):
            data["parameters"] = process_parameters(data["parameters"], "global.parameters")

        # 2. 处理paths
        if "paths" in data:
            data["paths"] = process_paths(data["paths"])

        # 3. 处理components
        if "components" in data:
            data["components"] = process_components(data["components"])

        # 4. 查找其他位置的example（不修改，只统计）
        if dry_run:
            print("\n[查找非参数位置的example字段]")
        find_other_examples(data)

        # 写入文件
        if not dry_run and output_file:
            with open(output_file, 'w', encoding='utf-8') as f:
                json.dump(data, f, ensure_ascii=False, indent=2)
            print(f"✓ 已保存到: {output_file}")
        elif dry_run:
            print(f"\n✓ 模拟执行完成，未实际修改文件")

        return stats

    except FileNotFoundError:
        print(f"✗ 错误: 文件 {input_file} 不存在")
        return {}
    except json.JSONDecodeError as e:
        print(f"✗ 错误: JSON解析失败 - {e}")
        return {}
    except Exception as e:
        print(f"✗ 错误: {e}")
        import traceback
        traceback.print_exc()
        return {}


def quick_fix_swagger(input_file: str, output_file: str = None) -> bool:
    """
    快速修复：修复所有参数中的example为x-example
    不询问，不交互，直接修复
    """
    if output_file is None:
        base, ext = os.path.splitext(input_file)
        output_file = f"{base}_fixed{ext}"

    try:
        with open(input_file, 'r', encoding='utf-8') as f:
            data = json.load(f)

        fixed_count = 0

        def fix_param_examples(params):
            """修复参数数组中的example字段"""
            nonlocal fixed_count
            if not isinstance(params, list):
                return params

            fixed = []
            for param in params:
                if not isinstance(param, dict):
                    fixed.append(param)
                    continue

                param_copy = param.copy()
                # 如果是参数（有"in"字段）且有example，则修复
                if "in" in param_copy and param_copy["in"] in ["path", "query", "header", "cookie"]:
                    if "example" in param_copy:
                        param_copy["x-example"] = param_copy.pop("example")
                        fixed_count += 1
                fixed.append(param_copy)

            return fixed

        # 处理根级别parameters
        if "parameters" in data and isinstance(data["parameters"], list):
            data["parameters"] = fix_param_examples(data["parameters"])

        # 处理paths
        if "paths" in data and isinstance(data["paths"], dict):
            for path, methods in data["paths"].items():
                if not isinstance(methods, dict):
                    continue
                for method, operation in methods.items():
                    if isinstance(operation, dict) and "parameters" in operation:
                        data["paths"][path][method]["parameters"] = fix_param_examples(operation["parameters"])

        # 处理components/parameters
        if "components" in data and isinstance(data["components"], dict):
            if "parameters" in data["components"] and isinstance(data["components"]["parameters"], dict):
                for param_name, param_def in data["components"]["parameters"].items():
                    if isinstance(param_def, dict) and "in" in param_def and "example" in param_def:
                        if param_def["in"] in ["path", "query", "header", "cookie"]:
                            data["components"]["parameters"][param_name]["x-example"] = data["components"]["parameters"][param_name].pop("example")
                            fixed_count += 1

        # 写入文件
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(data, f, ensure_ascii=False, indent=2)

        print(f"✓ 修复完成! 共修改了 {fixed_count} 个参数中的example字段")
        print(f"✓ 已保存到: {output_file}")
        return True

    except Exception as e:
        print(f"✗ 错误: {e}")
        return False


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("swagger文档example字段修复工具")
        print("=" * 50)
        print("用法:")
        print("  1. 自动修复所有参数中的example (快速模式):")
        print("     python fix_swagger.py <swagger文件路径>")
        print()
        print("  2. 模拟运行，只显示会修改的内容:")
        print("     python fix_swagger.py <swagger文件路径> --dry-run")
        print()
        print("  3. 快速修复并覆盖原文件:")
        print("     python fix_swagger.py <swagger文件路径> --quick")
        print()
        print("  4. 指定输出文件:")
        print("     python fix_swagger.py <输入文件> <输出文件>")
        print()
        print("示例:")
        print("  python fix_swagger.py swagger.json")
        print("  python fix_swagger.py swagger.json --dry-run")
        print("  python fix_swagger.py swagger.json --quick")
        print("  python fix_swagger.py swagger.json swagger_fixed.json")
        sys.exit(1)

    input_file = sys.argv[1]

    if not os.path.exists(input_file):
        print(f"错误: 文件不存在 - {input_file}")
        sys.exit(1)

    if len(sys.argv) > 2:
        if sys.argv[2] == "--dry-run":
            # 模拟运行
            print("模拟运行模式 (不会实际修改文件):")
            stats = fix_swagger_selective(input_file, dry_run=True)
            if stats:
                print("\n" + "="*50)
                print("统计信息:")
                print(f"  总参数数量: {stats.get('total_params', 0)}")
                print(f"  将修复的参数: {stats.get('params_fixed', 0)}")
                if stats.get('path_params', 0) > 0:
                    print(f"  - path参数: {stats.get('path_params', 0)}")
                if stats.get('query_params', 0) > 0:
                    print(f"  - query参数: {stats.get('query_params', 0)}")
                if stats.get('header_params', 0) > 0:
                    print(f"  - header参数: {stats.get('header_params', 0)}")
                if stats.get('cookie_params', 0) > 0:
                    print(f"  - cookie参数: {stats.get('cookie_params', 0)}")
                print(f"  保持原样的example: {stats.get('other_examples', 0)}")

        elif sys.argv[2] == "--quick":
            # 快速修复
            success = quick_fix_swagger(input_file, input_file)  # 覆盖原文件
            if success:
                print("\n✓ 快速修复完成，已覆盖原文件")

        elif sys.argv[2].endswith(".json"):
            # 指定输出文件
            output_file = sys.argv[2]
            stats = fix_swagger_selective(input_file, output_file)
            if stats:
                print("\n" + "="*50)
                print("修复完成! 统计信息:")
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
                print(f"  保持原样的example: {stats.get('other_examples', 0)}")
        else:
            print(f"未知参数: {sys.argv[2]}")
    else:
        # 默认行为：快速修复，不覆盖原文件
        base, ext = os.path.splitext(input_file)
        output_file = f"{base}_fixed{ext}"

        print(f"快速修复模式:")
        print(f"输入文件: {input_file}")
        print(f"输出文件: {output_file}")
        print()

        stats = fix_swagger_selective(input_file, output_file)
        if stats:
            print("\n" + "="*50)
            print("修复完成! 统计信息:")
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
            print(f"  保持原样的example: {stats.get('other_examples', 0)}")