# AWS-Todo-Project
API Gateway + DynamoDB + Lambda  

A simple API that provides CRUD (Create, Read, Update, Delete) operations for a simple TODO object. The API is implemented
using AWS Services and accessible via a public URL.

#### Assumptions:
1. The `id` (UUID) is passed in the request and not generated.
2. The response status for PUT create new To-Do object is set to `200` instead of `201` as it is specified in the API specification. 
3. The response status for GET To-Do object by id if `id` is not found is set to `400` instead of `404` as it is specified in the API specification.

### Endpoint specification
* Get all To-Do objects
    * `GET`: `/todo`
    * Response: A list of To-Do objects with status `200` or an empty list if no objects are stored.
        * CURL command example:
        *  `curl -X GET HOSTED-URL`

* Get To-Do object by `id`
    * `GET`: `/todo/{id}`
    * Response: The requested To-Do object with status `200` or status `400` if the object with provided `id` not found.
        * CURL command example:
        *  `curl -X GET \
           HOSTED-URL?id=54f56cdf-d24e-49de-9342-e0e9db1e3102`
      
* Create new To-Do object
    * `PUT`:`/todo`
    * Response: On successful creation, status `200` with success message or status `400` with error message when operation is not successful.
        * CURL command example:
        * `curl -X PUT \
          HOSTED-URL \
          -H 'content-type: application/json' \
          -d '{
          "id": "54f56cdf-d24e-49de-9342-e0e9db1e3102",
          "title": "API",
          "description": "Write todo specification"
          }'`

* Update existing To-Do object
    * `PUT`:`/todo/{id}`
    * Response: On successful update, status `200` with success message or status `400` with error message when operation is not successful.
        * CURL command example:
        *  `curl -X PUT \
           'HOSTED_URL?id=54f56cdf-d24e-49de-9342-e0e9db1e3102' \
           -H 'content-type: application/json' \
           -d '{
           "id": "54f56cdf-d24e-49de-9342-e0e9db1e3102",
           "title": "API",
           "description": "Write todo specification and add details"
           }'`

* Delete existing To-Do object
    * `DELETE`: `/todo/{id}`
    * Response: On successful delete, status `200` with success message or status `400` with error message when operation is not successful.
        * CURL command example:
        *  `curl -X DELETE \
           'HOSTED-URL?id=54f56cdf-d24e-49de-9342-e0e9db1e3102'`

_HOSTED-URL in the above curl command examples is the AWS URL where the application is hosted, they need to be replaced with the actual URL for the curl command to work_