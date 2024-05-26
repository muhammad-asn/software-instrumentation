### Checklist
- [x] Set global log level to debug.
- [x] All logs created must be written to both console and file. You must put the file in the logs directory with name app.log.
- [x] You may add as many log statements as you want, but you must have at least one log statement for each log level (debug, info, error). You can play around with the log level by setting the global log level to info or error to see the effect.
- [x] For every request received, a new request_id must be added to the log context. You may use google/uuid package to generate the request_id.
- [x] request_id must be printed in all logs created within the request-scope.
- [x] For each function, you must add the function name to the log context. For example, if there is a function named funcA, the log generated within that function must also contain func key, and funcA value as additional context for all logs generated within that function.