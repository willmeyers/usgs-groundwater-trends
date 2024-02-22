import os
import json
import ssl
from pathlib import Path
from datetime import datetime
from urllib import request


# Do not verify ssl
ssl._create_default_https_context = ssl._create_unverified_context


USGS_BASE_URL = "https://waterservices.usgs.gov/nwis/iv/?sites={}&parameterCd={}&startDT={}&endDT={}&siteStatus=all&format=json"


usgs_site_url = lambda site_no, parameter_cd, start, end: USGS_BASE_URL.format(site_no, parameter_cd, start, end)


def datetime_to_str(dt):
    return dt.strftime("%Y-%m-%dT%H:%M:%S.000-04:00")


def get_water_level_timeseries(site_no: str, parameter_cd: int, start: datetime, end: datetime) -> dict:
    url = usgs_site_url(site_no, parameter_cd, datetime_to_str(start), datetime_to_str(end))

    r = request.Request(url, headers={"user-agent": "development--xxx--2"})
    with request.urlopen(r) as resp:
        content = resp.read()
        encoding = resp.info().get_content_charset("utf-8")
        data = json.loads(content.decode(encoding))

    return data


with open(Path.cwd() / "data/usgs_groundwater_sites.json", "r") as fo:
    data = json.loads(fo.read())
    sites = data["value"]

    start = datetime(1980, 1, 1)
    end = datetime.today()

    site_files = os.listdir(Path.cwd() / "data/sites")
    site_files = [f.split(".")[0] for f in site_files if f.endswith(".json")]
    site_files = set(site_files)

    for idx, site in enumerate(sites):
        if site["SiteNumber"] in site_files:
            continue
        else:
            site_files.add(site["SiteNumber"])

        print(f"fetching {site['SiteNumber']}...\t{idx + 1} of {len(sites)}")
        data = get_water_level_timeseries(
            site_no=site["SiteNumber"],
            parameter_cd=site["ParameterCode"],
            start=start,
            end=end
        )

        with open(Path.cwd() / f"data/sites/{site['SiteNumber']}.json", "w+") as so:
            dump = json.dumps(data)
            so.write(dump)
            print(f"wrote {len(dump)} to {site['SiteNumber']}.json")
