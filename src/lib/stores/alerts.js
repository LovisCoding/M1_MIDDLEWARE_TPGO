import {writable} from "svelte/store";
import axios from "axios";

export const alertsError = writable("")

// Récupère les alertes pour une ressource spécifique
// Endpoint: GET /resources/{id}/alerts
export function getAlertsForResource(resourceId) {
    return axios.get(`/config_api/resources/${resourceId}/alerts`)
        .then((res) => {
            // CORRECTION ICI : Si res.data est null, on renvoie []
            return res.data || []; 
        }) 
        .catch((err) => {
            console.error(`Error alerts for res ${resourceId}`, err)
            return [] // En cas d'erreur HTTP, on renvoie aussi un tableau vide
        })
}

export function addAlert(resourceId, mailId) {
    alertsError.set("")
    return axios.post(`/config_api/resources/${resourceId}/alerts/${mailId}`)
}

export function removeAlert(resourceId, mailId) {
    alertsError.set("")
    return axios.delete(`/config_api/resources/${resourceId}/alerts/${mailId}`)
}