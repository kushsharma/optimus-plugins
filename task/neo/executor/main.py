import os
import requests
import json
from image_to_ascii import ImageToAscii

# path where secret will be mounted in docker container
SECRET_PATH = "/tmp/key.json"


def start():
    """
    Sends a http call to nasa api, parses response and prints potential hazardous
    objects in near earth orbit
    :return:
    """
    opt_config = fetch_config_from_optimus()

    # user configuration for date range
    range_start = opt_config["envs"]["RANGE_START"]
    range_end = opt_config["envs"]["RANGE_END"]

    # secret token required for NASA API being fetched from a file mounted as
    # volume by optimus executor
    with open(SECRET_PATH, "r") as f:
        api_key = json.load(f)['key']
    if api_key is None:
        raise Exception("invalid api token")

    # send the request for given date range
    r = requests.get(url="https://api.nasa.gov/neo/rest/v1/feed",
                     params={'start_date': range_start, 'end_date': range_end, 'api_key': api_key})

    # extracting data in json format
    print("for date range {} - {}".format(range_start, range_end))
    print_details(r.json())

    return


def fetch_config_from_optimus():
    """
    Fetch configuration inputs required to run this task for a single schedule day
    Configurations are fetched using optimus rest api
    :return:
    """
    # try printing os env to see what all we have for debugging
    # print(os.environ)

    # prepare request
    optimus_host = os.environ["OPTIMUS_HOSTNAME"]
    scheduled_at = os.environ["SCHEDULED_AT"]
    project_name = os.environ["PROJECT"]
    job_name = os.environ["JOB_NAME"]

    r = requests.post(url="http://{}/api/v1beta1/project/{}/job/{}/instance".format(optimus_host, project_name, job_name),
                      json={'scheduled_at': scheduled_at,
                            'instance_name': "neo",
                            'instance_type': "TASK"})
    instance = r.json()

    # print(instance)
    return instance["context"]


def print_details(jd):
    """
    Parse and calculate what we need from NASA endpoint response

    :param jd: json data fetched from NASA API
    :return:
    """
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


if __name__ == "__main__":
    start()


def print_image(range_start):
    """
    print image in ascii

    BROKEN AT THE MOMENT, SOMEONE FIX ME PLEASE
    :param range_start:
    :return:
    """

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
