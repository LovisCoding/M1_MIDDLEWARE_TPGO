<script>
    import {events, eventsError, getEvents} from "$lib/stores/events.js";
    import {onMount} from "svelte";
    import {SvelteMap} from "svelte/reactivity";
    import {getResource} from "$lib/stores/resources.js";
    import Calendar from '@event-calendar/core';
    import DayGrid from '@event-calendar/day-grid';

    const pickColorsFrom = [
        '#F6DC00', '#FFFFFF', '#AB69AD', '#FAC000', '#E7A2FD', '#DE9800',
    ]
    let agendaColors = new SvelteMap()
    let eventsResources = new SvelteMap()
    let plugins = [DayGrid];

    let options = $state({
        view: 'dayGridMonth',
        events: [],
        editable: false,
        eventContent: (info) => {
            let desc = info.event.extendedProps.description || "";
            let loc = info.event.extendedProps.location || "";
            return {html: `<div style="${getEventColors(info.event)}">
                <span class="fs-5 me-2 fw-bold">${info.event.start.getHours()}:${String(info.event.start.getMinutes()).padStart(2, '0')}</span>
                <span class="fs-5" title="${desc}">${info.event.title} <br/> <i class="fs-6">${loc}</i></span>
            </div>`}
        },
        eventBackgroundColor : '#593196' // Le fond est violet
    });

    onMount(() => {
        getEvents()
    })

    function getEventColors(event){
        let colors = []
        let color = ""

        let rIds = event.extendedProps.resourceIds || [];
        
        // CORRECTION ICI : Texte en BLANC (#FFFFFF) pour contraster avec le fond violet
        if (rIds.length === 0) {
             return `color: #FFFFFF;`; 
        }

        rIds.forEach((r) => {
            if(!agendaColors.get(r)){
                agendaColors.set(r, pickColorsFrom[agendaColors.size % pickColorsFrom.length])
            }
            colors.push(agendaColors.get(r))
        })
   
        if(colors.length > 1){
            let percent = 100 / colors.length
            color = "color: transparent;background-clip: text;background-image: linear-gradient(to right,"
            colors.forEach((c) => {
                color += `${c} ${percent}%,`
            })
            color = color.substring(0, color.length - 1);
            color += ");"
        } else {
            color = `color: ${colors[0]};`
        }
        return color
    }

    events.subscribe((values) => {
        const eventsList = [];
        
        if (Array.isArray(values)) {
            values.forEach((value) => {
                let e = {
                    id: value.id,
                    start: new Date(value.start),
                    end: new Date(value.end),
                    title: value.name,
                    extendedProps: {
                        description: value.description,
                        location: value.location,
                        resourceIds: value.agendaIds 
                    }
                }
                eventsList.push(e);

                (value.agendaIds || []).forEach((rId) => {
                    if(!eventsResources.get(rId)){
                        getResource(rId)
                            .then((data) => {
                                eventsResources.set(rId, data)
                            })
                            .catch(()=>{})
                    }
                })
            });
        }
        options.events = eventsList;
    })
</script>

<div class="d-flex flex-column">
    <h1 class="text-center color-yellow mb-lg-5">Events</h1>
    {#if $eventsError}
        <div  class="text-center text-danger form-text fs-5">
            An error occurred : {$eventsError}
        </div>
    {/if}

    <div class="d-flex mb-4 justify-content-center flex-wrap">
        {#each agendaColors as [r, c]}
            <div class="d-flex me-3 mb-2 align-items-center">
                <span class="p-3 me-2 border rounded" style="{`background-color: ${c};`}"></span>
                <span>{eventsResources.get(r) ? eventsResources.get(r).name : r}</span>
            </div>
        {/each}
    </div>

    <Calendar {plugins} {options}/>
</div>