<script lang="ts">
    import { onMount } from "svelte";
    import * as d3 from "d3";
    import { mouse, selectedSite, cache } from "../store";
    import type { Mouse, Site, Datapoint, SiteCache } from "./types";

    const API_URL = process.env.API_URL;

    let fetchedSites: SiteCache = {};
    let m: Mouse = { x: 0, y: 0 };
    let site: Site | null = null;
    let datapoints: Datapoint[] = [];
    let widthBreak: number = Infinity;
    let heightBreak: number = Infinity;

    const updateChart = () => {
        const xs = d3
            .scaleLinear()
            .domain(d3.extent(datapoints, (d) => d.Year) as [number, number])
            .range([45, 370]);

        const ys = d3
            .scaleLinear()
            .domain([
                d3.min(datapoints, (d) => d.Value) || 0,
                d3.max(datapoints, (d) => d.Value) || 0,
            ])
            .range([60, 5]);

        const xaxis = d3.axisBottom(xs);
        xaxis.ticks(3);
        xaxis.tickFormat((d) => `${d}`);
        const yaxis = d3.axisLeft(ys);
        yaxis.ticks(3);
        yaxis.tickFormat((d) => `${d}m`);

        const line = d3
            .line<Datapoint>()
            .x((d) => xs(d.Year))
            .y((d) => ys(d.Value))
            .curve(d3.curveCatmullRom.alpha(0.2));

        const svg = d3.select("#chart");
        svg.selectAll("*").remove();
        svg.append("path")
            .datum(datapoints)
            .attr("fill", "none")
            .attr("stroke", "royalblue")
            .attr("stroke-width", 1)
            .attr("transform", `translate(0, 0)`)
            .attr("d", line);
        svg.append("g").attr("transform", `translate(0, 62.5)`).call(xaxis);
        svg.append("g").attr("transform", `translate(45, 0)`).call(yaxis);
    };

    cache.subscribe((value) => {
        fetchedSites = value;
    });

    mouse.subscribe((value) => {
        m = value;
    });

    selectedSite.subscribe(async (value) => {
        if (value) {
            site = value;
            if (fetchedSites[site.ID]) {
                datapoints = fetchedSites[site.ID];
            } else {
                const r = await fetch(
                    `${API_URL}/datapoints?site_id=${site.ID}`,
                );
                const payload = await r.json();
                datapoints = payload.Datapoints as Datapoint[];
                cache.set({
                    ...fetchedSites,
                    [site.ID]: datapoints,
                } as SiteCache);
            }
        }

        updateChart();
    });

    onMount(() => {
        widthBreak = document.body.getBoundingClientRect().width / 2;
        heightBreak = document.body.getBoundingClientRect().height / 2;
    });
</script>

<div
    id="chart-tooltip"
    role="tooltip"
    style={`position: absolute; top: ${m.y < heightBreak ? m.y - 32 : m.y - 132}px; left: ${m.x < widthBreak ? m.x + 32 : m.x - 432}px;`}
>
    <div id="chart-container">
        {#if site}
            <span id="chart-title">{site.SiteName}</span>
            <svg id="chart" width="400" height="80"> </svg>
        {:else}
            <span id="chart-title">Start by hovering over a site</span>
        {/if}
    </div>
</div>

<style>
    #chart-tooltip {
        border: 1px solid royalblue;
        background: white;
        padding: 8px;
    }

    #chart-container {
        display: flex;
        flex-direction: column;
        gap: 2px;
        place-items: center;
        max-width: 400px;
    }

    #chart-title {
        font-size: 0.8rem;
        font-family: monospace;
        font-weight: 600;
    }

    .x-axis .tick line {
        stroke: #666;
    }
</style>
