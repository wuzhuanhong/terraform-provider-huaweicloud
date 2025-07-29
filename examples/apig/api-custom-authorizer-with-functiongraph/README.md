# Register an API with Custom Authorizer and FunctionGraph

This example provides best practice code for using Terraform to create an API Gateway instance with a custom authorizer
and FunctionGraph backend on HuaweiCloud. The example includes the creation of a VPC, subnet,
security group, ECS (for Nginx), FunctionGraph function, API Gateway instance, custom authorizer, custom response,
and EIP binding.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `ecs_instance_name` - The name of the ECS instance
* `ecs_instance_admin_pass` - The admin password of the ECS instance
* `function_name` - The name of the FunctionGraph function
* `function_code` - The code content of the function
* `apig_instance_name` - The name of the API Gateway instance
* `enterprise_project_id` - The ID of the enterprise project to which the APIG instance
  belongs (Required for enterprise users)
* `custom_authorizer_name` - The name of the custom authorizer
* `group_name` - The name of the API Gateway group
* `response_name` - The name of the API Gateway response
* `channel_name` - The name of the API Gateway channel
* `api_name` - The name of the API
* `api_request_path` - The request path of the API
* `api_backend_params` - The backend parameters of the API (default: [])
  - `type` - The type of the backend parameter
  - `name` - The name of the backend parameter
  - `location` - The location of the backend parameter
  - `value` - The value of the backend parameter
* `api_web_path` - The web path of the API

#### Optional Variables

* `availability_zone` - The availability zone to which the APIG instance belongs (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `ecs_instance_image_id` - The image ID of the ECS instance (default: "")
* `ecs_instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `ecs_instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `ecs_instance_flavor_cpu_core_count` - The number of CPU cores (default: 2)
* `ecs_instance_flavor_memory_size` - The memory size in GB (default: 4)
* `ecs_instance_image_visibility` - The visibility of the ECS image (default: "public")
* `ecs_instance_image_os_type` - The OS type of the ECS instance (default: "Ubuntu")
* `function_memory_size` - The memory size (MB) for the function (default: 128)
* `function_runtime` - The runtime environment for the function (default: "Python3.9")
* `function_timeout` - The timeout (seconds) for the function (default: 3)
* `function_handler` - The handler of the function (default: "index.handler")
* `function_code_type` - The code type of the function (default: "inline")
* `function_app` - The application name for the function (default: "default")
* `custom_authorizer_function_version` - The version of the function (default: "latest")
* `custom_authorizer_type` - The type of the custom authorizer (default: "FRONTEND")
* `custom_authorizer_network_type` - The network type of the custom authorizer (default: "V1")
* `custom_authorizer_cache_age` - The cache age of the custom authorizer (default: 0)
* `custom_authorizer_is_body_send` - Whether to send body in the custom authorizer (default: false)
* `custom_authorizer_use_data` - The user data for backend access authorization (default: null)
* `custom_authorizer_identity` - The identity list for the custom authorizer (default: [])
  - `name` - The name of the identity
  - `location` - The location of the identity
  - `validation` - The validation of the identity (default: null)
* `response_rules` - The rules for the API Gateway response (default: [])
  - `error_type` - The error type of the response
  - `body` - The body of the response
  - `status_code` - The status code of the response (default: null)
  - `headers` - The headers of the response (default: [])
    + `key` - The key of the header
    + `value` - The value of the header
* `channel_port` - The port of the API Gateway channel (default: 80)
* `channel_balance_strategy` - The balance strategy of the channel (default: 1)
* `api_type` - The type of the API (default: "public")
* `api_request_protocol` - The request protocol of the API (default: "HTTP")
* `api_request_method` - The request method of the API (default: "GET")
* `api_matching` - The matching rule of the API (default: "NORMAL")
* `api_web_request_method` - The web request method (default: "GET")
* `api_web_request_protocol` - The web request protocol (default: "HTTP")
* `api_web_timeout` - The web timeout in ms (default: 5000)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables. Example:

  ```hcl
  vpc_name                = "your_vpc_name"
  subnet_name             = "your_subnet_name"
  security_group_name     = "your_security_group_name"
  ecs_instance_name       = "your_ecs_instance_name"
  ecs_instance_admin_pass = "your_ecs_password"
  api_name                = "your_api_name"
  function_name           = "your_function_name"
  function_code           =  "your_function_code"
  apig_instance_name      = "your_apig_instance_name"
  enterprise_project_id   = "your_enterprise_project_id"
  custom_authorizer_name  = "your_custom_authorizer_name"
  group_name              = "your_group_name"
  response_name           = "your_response_name"
  response_rules          = "your_response_rules"
  channel_name            = "your_channel_name"
  api_request_path        = "your_api_request_path"
  api_backend_params      = "your_api_backend_params"
  api_web_path            = "your_api_web_path"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control.
* All resources will be created in the specified region.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.72.0 |
