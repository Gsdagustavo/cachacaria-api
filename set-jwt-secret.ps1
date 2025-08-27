# Set the JWT_SECRET environment variable for the user profile
[System.Environment]::SetEnvironmentVariable('JWT_SECRET', 'secret', [System.EnvironmentVariableTarget]::User)

# Verify the change by getting the environment variable
$env:JWT_SECRET
