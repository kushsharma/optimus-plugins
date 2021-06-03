import os
import requests
import json
from image_to_ascii import ImageToAscii

API_KEY = ""
SECRET_PATH = "/tmp/key.json"


def start():
    opt_config = fetch_optimus()
    # range_start = os.environ["RANGE_START"]
    # range_end = os.environ["RANGE_END"]

    with open(SECRET_PATH, "r") as f:
        API_KEY = json.load(f)['key']
    range_start = opt_config["envs"]["RANGE_START"]
    range_end = opt_config["envs"]["RANGE_END"]

    r = requests.get(url="https://api.nasa.gov/neo/rest/v1/feed",
                     params={'start_date': range_start, 'end_date': range_end, 'api_key': API_KEY})

    # extracting data in json format
    print_details(r.json())

    # TODO
    # print_image(range_start)
    return


def print_details(jd):
    element_count = jd['element_count']
    potentially_hazardous = []
    for search_date in jd['near_earth_objects'].keys():
        for neo in jd['near_earth_objects'][search_date]:
            if neo["is_potentially_hazardous_asteroid"] is True:
                potentially_hazardous.append({
                    "name": neo["name"],
                    "estimated_diameter_km": neo["estimated_diameter"]["kilometers"]["estimated_diameter_max"],
                    "relative_velocity_kmh": neo["close_approach_data"][0]["relative_velocity"]["kilometers_per_hour"]
                })

    print("total tracking: {}\npotential hazardous: {}".format(element_count, len(potentially_hazardous)))
    for haz in potentially_hazardous:
        print("Name: {}\nEstimated Diameter: {} km\nRelative Velocity: {} km/h\n\n".format(
            haz["name"],
            haz["estimated_diameter_km"],
            haz["relative_velocity_kmh"]
        ))
    return


def print_image(range_start):
    range_parts = range_start.split("-")

    # get json for all images of this date
    r = requests.get(url="https://epic.gsfc.nasa.gov/api/enhanced/date/{}".format(range_start),
                     params={'api_key': API_KEY})
    jd = r.json()

    # get the first image
    imageName = jd[0]["image"] + ".png"
    f = open('earth.png', 'wb')
    f.write(
        requests.get(
            url="https://epic.gsfc.nasa.gov/archive/natural/{}/{}/{}/png/{}".format(range_parts[0], range_parts[1],
                                                                                    range_parts[2], imageName),
            params={'api_key': API_KEY}).content
    )
    f.close()

    # convert to ascii
    ImageToAscii(imagePath="earth.png", outputFile="earth.txt")

    with open("earth.txt", "r") as earthFD:
        print(earthFD.read())
    return


def fetch_optimus():
    # print(os.environ)
    optimus_host = os.environ["OPTIMUS_HOSTNAME"]
    scheduled_at = os.environ["SCHEDULED_AT"]
    project_name = os.environ["PROJECT"]
    job_name = os.environ["JOB_NAME"]
    r = requests.post(url="http://{}/api/v1/project/{}/job/{}/instance".format(optimus_host, project_name, job_name),
                      json={'scheduled_at': scheduled_at, 'instance_name': "neo", 'instance_type': "TASK"})
    instance = r.json()
    # print(instance)
    return instance["context"]


if __name__ == "__main__":
    start()
