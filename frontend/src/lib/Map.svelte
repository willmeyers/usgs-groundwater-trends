<script lang="ts">
    import { onMount } from "svelte";
    import * as d3 from "d3";
    import { selectedSite } from "../store";
    import type { Site } from "./types";
    import geojson from "../assets/gz_2010_us_040_00_20m.json";

    const API_URL = process.env.API_URL;

    let sites: Site[] = [];
    const width = 900;
    const height = 600;

    // @ts-ignore
    const projection = d3
        .geoAlbersUsa()
        .translate([width / 2, height / 2])
        .scale(1000);

    const path = d3.geoPath().projection(projection);

    const states = geojson.features as d3.GeoPermissibleObjects[];

    onMount(async () => {
        const r = await fetch(`${API_URL}/sites`);
        const payload = (await r.json()) as { Sites: Site[] };
        sites = payload.Sites;
        sites = sites.map((site) => ({
            ...site,
            Location: JSON.parse(String(site.Location)),
        }));
    });

    const tsColorScale = d3.scaleLinear([-30, 30], ["red", "blue"]);

    const getSiteX2Coordinate = (site: Site): number | null => {
        const coords = projection(site.Location.coordinates);
        if (coords) {
            return 8 + site.TSEstimateInMM + coords[0];
        }

        return null;
    };

    const getSiteY2Coordinate = (site: Site): number | null => {
        const coords = projection(site.Location.coordinates);
        const direction = Math.sign(site.TSEstimateInMM);
        if (coords) {
            return -1 * direction * (2 + site.TSEstimateInMM) + coords[1];
        }

        return null;
    };

    const setSelectedSiteOnHover = (site: Site | null) => {
        selectedSite.set(site);
    };
</script>

<svg id="map" viewBox="0 0 900 600">
    <g fill="#eee" stroke="dimgray" stroke-width={0.2}>
        {#each states as feature, i}
            <path d={path(feature)} />
        {/each}
    </g>
    <g>
        {#each sites as site, i}
            {#if projection(site.Location.coordinates)}
                <line
                    tabindex={0}
                    x1={projection(site.Location.coordinates)[0]}
                    y1={projection(site.Location.coordinates)[1]}
                    x2={getSiteX2Coordinate(site)}
                    y2={getSiteY2Coordinate(site)}
                    stroke={tsColorScale(site.MKZ)}
                    stroke-width={1.5}
                    on:mouseover={() => setSelectedSiteOnHover(site)}
                    on:focus={() => setSelectedSiteOnHover(site)}
                    on:mouseleave={() => setSelectedSiteOnHover(null)}
                    on:blur={() => setSelectedSiteOnHover(null)}
                    role="button"
                />
            {/if}
        {/each}
    </g>
</svg>

<style>
    svg line {
        cursor: pointer;
        transition: all 0.2s;
    }
</style>
