# Events present in the system

## Working Day
| Event                | Description                          | Emitter | Subscribers |
| -------------------- | ------------------------------------ | ------- | ----------- |
| NEW_PATIENT          | New patient in clinic                | lobby   | -           |
| PATIENT_GONE         | Patient was forced out of the queue  | bench   | nurse       |
| PATIENT_IN_QUEUE     | Patient sat down on the bench        | bench   | nurse       |
| APPOINTMENT_FINISHED | Doctor finished working with patient | doctor  | nurse       |
| APPOINTMENT_STARTED  | Doctor started working with patient  | doctor  | -           |
