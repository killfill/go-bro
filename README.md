
:-{

CREATE SERVICE
curl http://username:password@localhost:3000/v2/service_instances/service_id -d '{  "service_id":        "service-guid-here",
  "plan_id":           "plan-guid-here",
  "organization_guid": "org-guid-here",
  "space_guid":        "space-guid-here"
}' -X PUT -i


DELETE SERVICE
curl http://username:password@localhost:3000/v2/service_instances/service_id -X DELETE -i

BIND
curl http://username:password@localhost:3000/v2/service_instances/service_id/service_bindings/bind_id -d '{
  "plan_id": "plan-guid-here",
  "service_id": "service-guid-here",
  "app_guid": "app-guid-here"
}' -X PUT -i


UNBIND
curl 'http://username:password@localhost:3000/v2/service_instances/service_id/service_bindings/bind_id?service_id=service-id-here&plan_id=plan-id-here' -X DELETE -i
