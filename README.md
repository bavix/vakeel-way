# vakeel-way

Vakeel Way is a service event storage. The service is designed to store information about service events that are sent to it. The service provides the ability to update a list of UUIDs. The list of UUIDs is sent as a stream of UpdateRequest messages. Each UpdateRequest message contains a list of UUIDs that need to be updated. These UUIDs are used to uniquely identify the request and can be used to track the request throughout the system.

The service is used to test that the gRPC service is working correctly. This is done by sending a stream of UpdateRequest messages and receiving a single UpdateResponse message in return. The UpdateResponse message is an empty message that indicates that the update operation was successful.

It is also used to mark services as working for some time. If services stop sending information about themselves, then they do not work and it is necessary to notify monitoring and create an incident.
