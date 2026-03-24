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

Write-Host "1) Create user..."
$suffix = Get-Date -Format "yyyyMMddHHmmssfff"
$createUserBody = @{ email = "smoke.user.$suffix@example.com"; password = "secret123" } | ConvertTo-Json -Compress
$createdUser = Invoke-RestMethod -Method POST -Uri "$BaseUrl/users" -ContentType "application/json" -Body $createUserBody
Assert-True ($null -ne $createdUser.id) "Create user failed: response has no id"
$userId = [int]$createdUser.id
Write-Host "   Created user id = $userId"

Write-Host "2) Create task..."
$createBody = @{ task = "Buy milk"; is_done = $false; user_id = $userId } | ConvertTo-Json -Compress
$created = Invoke-RestMethod -Method POST -Uri "$BaseUrl/tasks" -ContentType "application/json" -Body $createBody
Assert-True ($null -ne $created.id) "Create failed: response has no id"
$taskId = [int]$created.id
Write-Host "   Created id = $taskId"

Write-Host "3) List tasks..."
$list1 = Invoke-RestMethod -Method GET -Uri "$BaseUrl/tasks"
if ($list1 -isnot [System.Array]) {
    $list1 = @($list1)
}
Assert-True ($list1.Count -ge 1) "List failed: no tasks returned after create"
Write-Host "   Tasks count = $($list1.Count)"

Write-Host "4) List tasks by user..."
$userTasks = Invoke-RestMethod -Method GET -Uri "$BaseUrl/users/$userId/tasks"
if ($userTasks -isnot [System.Array]) {
    $userTasks = @($userTasks)
}
Assert-True ($userTasks.Count -ge 1) "List by user failed: no tasks returned for user"
Write-Host "   User tasks count = $($userTasks.Count)"

Write-Host "5) Patch task..."
$patchBody = @{ is_done = $true } | ConvertTo-Json -Compress
$patched = Invoke-RestMethod -Method PATCH -Uri "$BaseUrl/tasks/$taskId" -ContentType "application/json" -Body $patchBody
Assert-True ($patched.is_done -eq $true) "Patch failed: is_done is not true"
Write-Host "   Patched is_done = $($patched.is_done)"

Write-Host "6) Delete task..."
Invoke-RestMethod -Method DELETE -Uri "$BaseUrl/tasks/$taskId" | Out-Null
Write-Host "   Deleted id = $taskId"

Write-Host "7) Delete user..."
Invoke-RestMethod -Method DELETE -Uri "$BaseUrl/users/$userId" | Out-Null
Write-Host "   Deleted user id = $userId"

Write-Host "8) Final list..."
$list2 = Invoke-RestMethod -Method GET -Uri "$BaseUrl/tasks"
if ($list2 -isnot [System.Array]) {
    $list2 = @($list2)
}
Write-Host "   Final tasks count = $($list2.Count)"

Write-Host "Smoke test passed."
