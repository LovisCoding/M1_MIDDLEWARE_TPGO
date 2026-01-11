import {writable} from "svelte/store";
import axios from "axios";

export const events = writable([])
export const eventsError = writable("")

const apiBaseUrl = "/timetable_api/events"

/*
***** BACKEND DATA FORMAT (Swagger)
* {
* "id": "...",
* "summary": "TD Entrepôt de données - G1",
* "location": "IS_A104",
* "description": "...",
* "startTime": "...",
* "endTime": "...",
* "resourceIds": [...]
* }
*/

export function getEvents() {
    eventsError.set("")
    axios.get(`${apiBaseUrl}`)
        .then((res) => {
            // Sécurité : on vérifie que res.data est bien un tableau
            if (Array.isArray(res.data)) {
                const mappedEvents = res.data.map(e => ({
                    id: e.id,
                    name: e.summary, 
                    description: e.description,
                    location: e.location,
                    start: e.startTime,
                    end: e.endTime,
                    // On s'assure que agendaIds est un tableau
                    agendaIds: Array.isArray(e.resourceIds) ? e.resourceIds : [] 
                }))
                events.set(mappedEvents)
            } else {
                console.error("Format de données inattendu", res.data);
                events.set([]);
            }
        })
        .catch((err) => {
            console.log("An error has occurred while retrieving events")
            console.log(err)
            if(err.response?.data?.message){
                eventsError.set(JSON.stringify(err.response.data.message))
            } else {
                eventsError.set(JSON.stringify(err))
            }
        })
}