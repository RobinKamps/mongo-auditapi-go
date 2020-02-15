# MongoDB Audit Trail API
This project complements the MongoDB Change Stream Watcher, providing a REST API that can be used to fetch field-level audit trails 
from audit records created by the change stream watcher.

### Structure of Audit Records in the Audit Database
The following is an example of an audit record that is created by the watcher:
```json
{ 
    "_id" : {
        "_data" : "825E446661000000012B022C0100296E5A1004EC1E76078DCE4C489A2BFE17218EC79F46645F696400645C5D85C62FEF357A165CCABF0004"
    }, 
    "collection" : "streamtest", 
    "database" : "test", 
    "documentKey" : "5c5d85c62fef357a165ccabf", 
    "fullDocument" : null, 
    "operationType" : "update", 
    "timestamp" : "2020-02-14T15:03:27Z", 
    "updateDescription" : {
        "updatedFields" : {
            "lineItems.0.procedures.0.procedureModCodes.0" : "332"
        }, 
        "removedFields" : [

        ]
    }, 
    "user" : "tcadmin"
}
```
Each audit record captures changes to a certain record in the MongoDB collection that is being watched by the change watcher.
The "documentKey" identifies the audited record. The audit record contains details of the change, such as what was the change operation
(insert/update/delete, etc.), when was the change made, by whom, what fields changed, and what were the values of the update fields following 
the change. For "insert" operations, the entire record is included in the "fullDocument" field.

### API Endpoint
The API has a single endpoint that can be used to fetch the audit trail of a specific field of a specific record.
> GET /auditrecords/{documentKey}/{fieldPath}

**documentKey**: The ID of the record which contains the field for which the audit trail needs to be fetched. <br>  
**fieldPath**: The fully qualified JSON path of the field whose audit trail needs to be fetched. <br>  

### Sample Response
The following is a sample response for the request GET /auditrecords/5c5d85c62fef357a165ccabf/lineItems.0.procedures.0.procedureModCodes.0

```json
[
    {
        "fieldId": "lineItems.0.procedures.0.procedureModCodes.0",
        "fieldValue": "332",
        "updatedBy": "jdoe",
        "updatedAt": "2020-02-12T15:56:01-05:00"
    },
    {
        "fieldId": "lineItems.0.procedures.0.procedureModCodes.0",
        "fieldValue": "453",
        "updatedBy": "ssmith",
        "updatedAt": "2020-02-12T15:32:33-05:00"
    },
    {
        "fieldId": "lineItems.0.procedures.0.procedureModCodes.0",
        "fieldValue": "444",
        "updatedBy": "mwilson",
        "updatedAt": "2020-02-12T12:36:40-05:00"
    },
    {
        "fieldId": "lineItems.0.procedures.0.procedureModCodes.0",
        "fieldValue": "303",
        "updatedBy": "tparker",
        "updatedAt": "2020-01-30T12:06:52-05:00"
    }
]

```
## Running the program
**Note**: The current version of the program has been tested with Golang 1.13.5. It utilizes Go modules. 
1. Build the API program using "go build -o mongo-auditapi cmd/main.go"
2. Copy the config file config/config.json to an appropriate location - modify per your environment.
3. Set the location of config.json via environment variable CONFIG_FILE.
4. Execute the API program "./mongo-auditapi".


