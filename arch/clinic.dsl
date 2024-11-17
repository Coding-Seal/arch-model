workspace "Arch-Model" "Assignment in polytech" {

    !identifiers hierarchical

    model {
        u = person "Пациент"
        ss = softwareSystem "Поликлиника" {
            lobby = container "Регистратура"{
                tags "Many"
            }
            bench = container "Скамека" {
                tags "Storage"
            }
            taskExtractor = container "Медсестра"
            
            taskExecutor = container "Врач"{
                tags "Many"
            }
            eventManager = container "Блок обработки событий"
        }
        u -> ss "Приходит на прием"
        u -> ss.lobby "Идет получать талон"

        ss.lobby -> ss.eventManager "Send PATIENT_CAME event"
        ss.lobby -> ss.bench "Человек с талоном садится на скамейку"

        ss.bench -> ss.eventManager "События PATIENT_GONE OR/AND PATIENT_IN_QUEUE "

        ss.taskExtractor -> ss.bench "Зовет пациента"
        ss.taskExtractor -> ss.taskExecutor "Отправляет пациента к врачу"
        ss.taskExtractor -> ss.eventManager "Событие CALLED_PATIENT"

        ss.taskExecutor -> ss.eventManager "Событие TREATMENT_DONE"

        ss.eventManager -> ss.taskExtractor "Уведомить медсестру, что врач свободен" 
        ss.eventManager -> ss.taskExtractor "Уведомить медсестру, что пришел новый пациент"       
    }

    views {
        systemContext ss "Hospital" {
            include *
            autolayout
        }
        container ss "Components" {
            include *
            autolayout lr
        }
        dynamic ss {
            title "Browse top 20 books feature"
   
            autoLayout lr
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