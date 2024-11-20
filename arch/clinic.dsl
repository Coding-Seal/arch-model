workspace "Arch-Model" "Assignment in polytech" {

    !identifiers hierarchical

    model {
        u = person "Patient"
        ss = softwareSystem "Clinic" {
            lobby = container "Lobby"{
                tags "Many"
            }
            bench = container "Bench" {
                tags "Storage"
            }
            nurse = container "Nurse"
            
            doctor = container "Doctor"{
                tags "Many"
            }
            eventManager = container "Event Manager"
        }
        u -> ss "Comes to clinic"
        u -> ss.lobby "Comes to lobby"

        ss.lobby -> ss.eventManager "[EVENT] NEW_PATIENT"
        ss.lobby -> ss.bench "Tries to sit down"

        ss.bench -> ss.eventManager "[EVENT] PATIENT_IN_QUEUE"
        ss.bench -> ss.eventManager "[EVENT] PATIENT_GONE"

        ss.nurse -> ss.bench "Call patient"
        ss.nurse -> ss.doctor "Send patient to the doctor"


        ss.doctor -> ss.eventManager "[EVENT] APPOINTMENT_STARTED"
        ss.doctor -> ss.eventManager "[EVENT] APPOINTMENT_FINISHED"

        ss.eventManager -> ss.nurse "[EVENT] APPOINTMENT_FINISHED" 
        ss.eventManager -> ss.nurse "[EVENT] PATIENT_IN_QUEUE"
        ss.eventManager -> ss.nurse "[EVENT] PATIENT_GONE"        
    }

    views {
        systemContext ss "Hospital" {
            include *
            autolayout
        }
        container ss "Components" {
            include *
            autolayout
        }

        styles { 
            element "Person" {
                shape Person
                color white
                background SteelBlue
            }
            element "Element" {
                color white
                background SteelBlue
                shape RoundedBox
            }
            element "Storage"{
                shape Cylinder
            }
            element "Storage"{
                shape Cylinder
            }
        }
    }

}