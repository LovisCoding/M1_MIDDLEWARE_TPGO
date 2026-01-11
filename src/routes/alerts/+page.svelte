<script>
    import {getAlertsForResource, addAlert, removeAlert, alertsError} from "$lib/stores/alerts.js";
    import {resources, getResources} from "$lib/stores/resources.js";
    import {mails, getMails, postMail} from "$lib/stores/mails.js";
    import DeleteModal from "../../components/delete-modal.svelte"
    import {onMount} from "svelte";

    let newEmail = $state("")
    let selectedResourceId = $state("")
    let displayedAlerts = $state([])
    let deleteModal = $state();

    onMount(async () => {
        await getResources();
        await getMails();
    })

    // Recharge la liste quand les ressources changent
    $effect(() => {
        if ($resources.length > 0) loadAllAlerts();
    });

    async function loadAllAlerts() {
        let list = [];
        for (const res of $resources) {
            const resAlerts = await getAlertsForResource(res.id); // Renvoie [{mail_id, ...}]
            resAlerts.forEach(a => {
                // Le backend renvoie 'mail_id' (snake_case)
                const mId = a.mail_id || a.mailId; 
                const mailObj = $mails.find(m => m.id === mId);
                if (mailObj) {
                    list.push({
                        uniqueId: `${res.id}-${mailObj.id}`, // ID composite pour la suppression
                        resId: res.id,
                        resName: res.name,
                        mailId: mailObj.id,
                        email: mailObj.mail
                    });
                }
            });
        }
        displayedAlerts = list;
    }

    async function handleAdd(e) {
        e.preventDefault();
        if (!newEmail || !selectedResourceId) return;

        // 1. Trouve ou Crée le mail
        let mailObj = $mails.find(m => m.mail === newEmail);
        if (!mailObj) {
            mailObj = await postMail(newEmail);
        }

        // 2. Ajoute le lien
        if (mailObj && mailObj.id) {
            await addAlert(selectedResourceId, mailObj.id);
            newEmail = "";
            selectedResourceId = "";
            loadAllAlerts();
        }
    }

    function handleDelete(uniqueId) {
        const [rId, mId] = uniqueId.split("-");
        removeAlert(rId, mId).then(() => loadAllAlerts());
    }
</script>

<div class="d-flex flex-column">
    <h1 class="text-center color-yellow mb-lg-5">Alerts</h1>
    
    <div class="d-flex flex-column mb-5">
        <h5 class="text-center fw-semibold mb-2">Link Email to Resource : </h5>
        <form class="d-flex justify-content-center align-items-end" onsubmit={handleAdd}>
            <div class="d-flex flex-column me-4">
                <label for="res" class="fs-5 mb-1">Resource</label>
                <select id="res" class="form-select form-select-sm" bind:value={selectedResourceId} required>
                    <option value="" disabled selected>Select resource</option>
                    {#each $resources as r} <option value={r.id}>{r.name} ({r.type})</option> {/each}
                </select>
            </div>
            <div class="d-flex flex-column me-4">
                <label for="mail" class="fs-5 mb-1">Email</label>
                <input id="mail" type="email" class="form-control form-control-sm" placeholder="user@uca.fr" bind:value={newEmail} required/>
            </div>
            <input class="btn btn-yellow" type="submit" value="Link"/>
        </form>
    </div>

    {#if $alertsError} <div class="text-center text-danger">{$alertsError}</div> {/if}

    <table class="table w-50 align-self-center">
        <thead>
            <tr><th class="color-yellow">Resource</th><th class="color-yellow">Email</th><th>Action</th></tr>
        </thead>
        <tbody>
        {#each displayedAlerts as item (item.uniqueId)}
            <tr>
                <td>{item.resName}</td>
                <td>{item.email}</td>
                <td>
                    <button class="bg-transparent border-0" onclick={() => deleteModal.show(item.uniqueId, `${item.email} on ${item.resName}`)}>
                        <span class="fa fa-trash-can text-danger"></span>
                    </button>
                </td>
            </tr>
        {:else}
            <tr><td colspan="3" class="text-center">No alerts configured.</td></tr>
        {/each}
        </tbody>
    </table>
    <DeleteModal bind:this={deleteModal} uniqueName="link" deleteFunction={handleDelete} />
</div>