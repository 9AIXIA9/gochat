# 统计Go语言文件行数（排除_gen后缀文件，可选排除注释和空行）
param(
    [string]$Path = ".",
    [switch]$Verbose = $false,
    [switch]$ExcludeComments = $false,
    [switch]$ExcludeEmptyLines = $false
)

function Count-CodeLines {
    param(
        [string]$filePath,
        [bool]$excludeComments,
        [bool]$excludeEmptyLines
    )

    $content = Get-Content $filePath
    $totalLines = 0
    $inBlockComment = $false

    foreach ($line in $content) {
        $trimmedLine = $line.Trim()

        # 跳过空行
        if ($excludeEmptyLines -and [string]::IsNullOrWhiteSpace($trimmedLine)) {
            continue
        }

        # 处理块注释
        if ($excludeComments) {
            # 跳过块注释内的行
            if ($inBlockComment) {
                # 检查块注释是否结束
                if ($trimmedLine -match "\*/") {
                    $inBlockComment = $false

                    # 检查块注释结束符后面是否有代码
                    $afterComment = $trimmedLine -replace "^.*?\*/", ""
                    if (-not [string]::IsNullOrWhiteSpace($afterComment) -and
                        -not ($afterComment.Trim().StartsWith("//"))) {
                        $totalLines++
                    }
                }
                continue
            }

            # 检查块注释开始
            if ($trimmedLine -match "/\*") {
                $inBlockComment = $true

                # 检查块注释开始符前面是否有代码
                $beforeComment = $trimmedLine -replace "/\*.*", ""
                if (-not [string]::IsNullOrWhiteSpace($beforeComment) -and
                    -not ($beforeComment.Trim().EndsWith("//"))) {
                    $totalLines++
                }

                # 检查块注释是否在同一行结束
                if ($trimmedLine -match "\*/") {
                    $inBlockComment = $false

                    # 检查块注释结束符后面是否有代码
                    $afterComment = $trimmedLine -replace "^.*?\*/", ""
                    if (-not [string]::IsNullOrWhiteSpace($afterComment) -and
                        -not ($afterComment.Trim().StartsWith("//"))) {
                        $totalLines++
                    }
                }
                continue
            }

            # 跳过单行注释
            if ($trimmedLine.StartsWith("//")) {
                continue
            }

            # 跳过包含行末注释的行（如果整行只有注释）
            if ($trimmedLine -match "//" -and ($trimmedLine -replace "//.*", "").Trim() -eq "") {
                continue
            }
        }

        $totalLines++
    }

    return $totalLines
}

function Get-GoFileLineCount {
    param([string]$Directory)

    $totalLines = 0
    $codeLines = 0
    $fileCount = 0

    # 获取所有.go文件，排除_gen后缀的文件
    $goFiles = Get-ChildItem -Path $Directory -Recurse -Filter "*.go" |
               Where-Object { $_.Name -notlike "*_gen.go" }

    foreach ($file in $goFiles) {
        try {
            # 统计总行数
            $allLines = (Get-Content $file.FullName | Measure-Object -Line).Lines

            # 统计代码行数（排除注释和空行）
            if ($ExcludeComments -or $ExcludeEmptyLines) {
                $codeLineCount = Count-CodeLines -filePath $file.FullName `
                                                -excludeComments $ExcludeComments `
                                                -excludeEmptyLines $ExcludeEmptyLines
            } else {
                $codeLineCount = $allLines
            }

            $totalLines += $allLines
            $codeLines += $codeLineCount
            $fileCount++

            if ($Verbose) {
                $info = "文件: $($file.Name) - 总行数: $allLines"
                if ($ExcludeComments -or $ExcludeEmptyLines) {
                    $info += " - 代码行数: $codeLineCount"
                }
                Write-Host $info
            }
        }
        catch {
            Write-Warning "无法读取文件: $($file.FullName)"
        }
    }

    return @{
        TotalLines = $totalLines
        CodeLines = $codeLines
        FileCount = $fileCount
    }
}

Write-Host "正在统计Go语言文件行数..." -ForegroundColor Green
Write-Host "排除规则: *_gen.go 文件" -ForegroundColor Yellow
if ($ExcludeComments) {
    Write-Host "排除注释: 是" -ForegroundColor Yellow
}
if ($ExcludeEmptyLines) {
    Write-Host "排除空行: 是" -ForegroundColor Yellow
}
Write-Host "搜索路径: $Path" -ForegroundColor Yellow
Write-Host ("-" * 50)

$result = Get-GoFileLineCount -Directory $Path

Write-Host ("-" * 50)
Write-Host "统计结果:" -ForegroundColor Green
Write-Host "文件数量: $($result.FileCount)" -ForegroundColor Cyan
Write-Host "总行数: $($result.TotalLines)" -ForegroundColor Cyan
if ($ExcludeComments -or $ExcludeEmptyLines) {
    Write-Host "代码行数: $($result.CodeLines)" -ForegroundColor Green
    $percentage = [math]::Round(($result.CodeLines / $result.TotalLines * 100), 2)
    Write-Host "代码比例: $percentage%" -ForegroundColor Green
}

# 显示详细信息（如果文件数量较少）
if ($result.FileCount -le 20 -and !$Verbose) {
    Write-Host "`n提示: 使用 -Verbose 参数查看每个文件的详细行数" -ForegroundColor Gray
}