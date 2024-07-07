member_id='xxxx-yyyy-zzzz' # member id (local id for user, origin id for group)
role_name='Administrator'  # role name (one of: 'Administrator', 'User', 'Reader')
environment_id='1234'      # environment id

terraform import azuredevopsext_environment_security_access "${member_id},${role_name}@${environment_id}"
