<script lang="ts">
    import { mouse } from "./store";
    import Chart from "./lib/Chart.svelte";
    import Map from "./lib/Map.svelte";

    const handleMouseMove = (event: any) => {
        const position = { x: event.clientX, y: event.clientY };
        mouse.set(position);
    };
</script>

<main on:mousemove={handleMouseMove}>
    <header>
        <h1>USGS Groundwater Trends</h1>
    </header>
    <Chart />
    <Map />
    <footer>
        <a href="https://willmeye.rs">made by will meyers</a>
        <small
            >Trends were determined with Mann-Kendall and Theil-Sen median
            regressions. Data sourced via USGS.</small
        >
    </footer>
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        width: 100vw;
        height: 100vh;
    }

    header {
        display: flex;
        justify-content: center;
        align-items: center;
    }

    header h1 {
        margin: 4px;
        font-size: 1rem;
    }

    footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
</style>
