pipeline_id='1234'    # pipeline id
environment_id='1234' # environment id

terraform import azuredevopsext_environment_security_access "${pipeline_id}@${environment_id}'
