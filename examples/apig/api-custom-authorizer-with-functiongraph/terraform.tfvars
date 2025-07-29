vpc_name                = "tf_test_apig_api_auth"
subnet_name             = "tf_test_apig_api_auth"
security_group_name     = "tf_test_apig_api_auth"
ecs_instance_name       = "tf_test_apig_api_auth"
ecs_instance_admin_pass = "YourECSInstancePassword!"
function_name           = "tf_test_function"
function_code           = <<EOF
# -*- coding:utf-8 -*-
import json
def handler(event, context):
    if event["headers"].get("x-user-auth")=='cXpsdzQyVW9Xa1NVTX==':
        return {
            'statusCode': 200,
            'body': json.dumps({
                "status":"allow",
                "context":{
                    "user_name":"user1"
                }
            })
        }
    else:
        return {
            'statusCode': 200,
            'body': json.dumps({
                "status":"deny",
                "context":{
                    "code":"1001",
                    "message":"incorrect username or password"
                }
            })
        }
EOF
apig_instance_name      = "tf_test_apig_instance"
enterprise_project_id   = "0"
custom_authorizer_name  = "tf_test_custom_authorizer"
group_name              = "tf_test_apig_group"
response_name           = "tf_test_apig_response"
response_rules = [{
  error_type  = "AUTHORIZER_FAILURE"
  body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
  status_code = 401
}]

channel_name     = "tf_test_apig_channel"
api_request_path = "/backend/users"
api_backend_params = [{
  type     = "SYSTEM"
  name     = "X-User-Auth"
  location = "HEADER"
  value    = "user_name"
}]
api_web_path = "/backend/users"
api_name     = "tf_test_apig_api_auth"


region_name = "cn-north-4"
access_key  = "IUFN87TS8WZPWJRX9QLH"
secret_key  = "xoZ4asqmAGrchDsGwkE3w0DoS5tRaugtp1uigfg3"
