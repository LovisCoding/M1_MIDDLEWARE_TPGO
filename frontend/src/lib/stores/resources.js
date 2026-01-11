import {writable} from "svelte/store";
import axios from "axios";

export const resources = writable([])
export const resourcesError = writable("")

const apiBaseUrl = "/config_api/resources"

export function getResources() {
    resourcesError.set("")
    axios.get(`${apiBaseUrl}`)
        .then((res) => resources.set(res.data))
        .catch((err) => handleError(err, resourcesError))
}

export function postResource(name, type) {
    resourcesError.set("")
    return axios.post(`${apiBaseUrl}`, { name, type })
        .then(() => getResources())
        .catch((err) => handleError(err, resourcesError))
}

export function putResource(id, name, type) {
    resourcesError.set("")
    return axios.put(`${apiBaseUrl}/${id}`, { name, type })
        .then(() => getResources())
        .catch((err) => handleError(err, resourcesError))
}

export function deleteResource(id) {
    resourcesError.set("")
    return axios.delete(`${apiBaseUrl}/${id}`)
        .then(() => getResources())
        .catch((err) => handleError(err, resourcesError))
}

export function getResource(id) {
    return axios.get(`${apiBaseUrl}/${id}`).then(res => res.data);
}

function handleError(err, store) {
    console.error(err);
    if(err.response?.data?.message){
        store.set(JSON.stringify(err.response.data.message))
    } else {
        store.set(JSON.stringify(err.message))
    }
}