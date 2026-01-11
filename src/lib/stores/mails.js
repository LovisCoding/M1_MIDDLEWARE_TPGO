import {writable} from "svelte/store";
import axios from "axios";

export const mails = writable([])
export const mailsError = writable("")

const apiBaseUrl = "/config_api/mails"

/*
***** DATA FORMAT EXAMPLE
* {
* "id": 1,
* "mail": "example@etu.uca.fr"
* }
*/

export function getMails() {
    mailsError.set("")
    return axios.get(`${apiBaseUrl}`)
        .then((res) => {
            mails.set(res.data)
            return Promise.resolve(res.data)
        })
        .catch((err) => {
            console.log("An error has occurred while retrieving mails")
            if(err.response?.data?.message){
                mailsError.set(JSON.stringify(err.response.data.message))
            } else {
                mailsError.set(JSON.stringify(err))
            }
        })
}

export function postMail(email) {
    mailsError.set("")
    // Le backend attend un objet avec le champ "mail" (inferred from schema)
    // Mais vérifiez si votre contrôleur mails attend "mail" ou "email"
    // Selon le schéma DB: "mail TEXT NOT NULL UNIQUE", on suppose JSON { "mail": "..." }
    return axios.post(`${apiBaseUrl}`, { mail: email })
        .then((res) => {
            getMails()
            return Promise.resolve(res.data)
        })
        .catch((err) => {
            console.log("An error as occurred while posting mail")
            if(err.response?.data?.message){
                mailsError.set(JSON.stringify(err.response.data.message))
            } else {
                mailsError.set(JSON.stringify(err))
            }
            return Promise.reject(err)
        })
}

export function deleteMail(id) {
    mailsError.set("")
    return axios.delete(`${apiBaseUrl}/${id}`)
        .then(() => {
            getMails()
            return Promise.resolve()
        })
        .catch((err) => {
            console.log("An error as occurred while deleting mail")
            if(err.response?.data?.message){
                mailsError.set(JSON.stringify(err.response.data.message))
            } else {
                mailsError.set(JSON.stringify(err))
            }
            return Promise.reject(err)
        })
}