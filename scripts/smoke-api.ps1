param(
    [string]$BaseUrl = "http://localhost:8080"
)

$ErrorActionPreference = "Stop"

function Assert-True {
    param(
        [bool]$Condition,
        [string]$Message
    )
    if (-not $Condition) {
        throw $Message
    }
}

Write-Host "Base URL: $BaseUrl"

Write-Host "1) Create task..."
$createBody = @{ task = "Buy milk"; is_done = $false } | ConvertTo-Json -Compress
$created = Invoke-RestMethod -Method POST -Uri "$BaseUrl/tasks" -ContentType "application/json" -Body $createBody
Assert-True ($null -ne $created.id) "Create failed: response has no id"
$taskId = [int]$created.id
Write-Host "   Created id = $taskId"

Write-Host "2) List tasks..."
$list1 = Invoke-RestMethod -Method GET -Uri "$BaseUrl/tasks"
if ($list1 -isnot [System.Array]) {
    $list1 = @($list1)
}
Assert-True ($list1.Count -ge 1) "List failed: no tasks returned after create"
Write-Host "   Tasks count = $($list1.Count)"

Write-Host "3) Patch task..."
$patchBody = @{ is_done = $true } | ConvertTo-Json -Compress
$patched = Invoke-RestMethod -Method PATCH -Uri "$BaseUrl/tasks/$taskId" -ContentType "application/json" -Body $patchBody
Assert-True ($patched.is_done -eq $true) "Patch failed: is_done is not true"
Write-Host "   Patched is_done = $($patched.is_done)"

Write-Host "4) Delete task..."
Invoke-RestMethod -Method DELETE -Uri "$BaseUrl/tasks/$taskId" | Out-Null
Write-Host "   Deleted id = $taskId"

Write-Host "5) Final list..."
$list2 = Invoke-RestMethod -Method GET -Uri "$BaseUrl/tasks"
if ($list2 -isnot [System.Array]) {
    $list2 = @($list2)
}
Write-Host "   Final tasks count = $($list2.Count)"

Write-Host "Smoke test passed."
