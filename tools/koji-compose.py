#!/usr/bin/python3
import json
import sys
import time

import requests

DISTRO_REPOS = {
    "fedora-31": [
        {
            "baseurl": "http://download.fedoraproject.org/pub/fedora/linux/releases/31/Everything/x86_64/os/",
            "gpgkey": "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmQINBFxq3QMBEADUhGfCfP1ijiggBuVbR/pBDSWMC3TWbfC8pt7fhZkYrilzfWUM\nfTsikPymSriScONXP6DNyZ5r7tgrIVdVrJvRIqIFRO0mufp9HyfWKDO//Ctyp7OQ\nzYw6NVthO/aWpyFfJpj6s4iZsYGqf9gByV8brBB8v8jEsCtVOj1BU3bMbLkMsRI9\n+WiLjDYyvopqNBQuIe8ogxSxpYdbUz6+jxzfvhRoBzWdjITd//Gjd90kkrBOMWkO\nLTqO133OD1WMT08G5NuQ4KhjYsVvSbBpfdkTcNuP8gBP9LxCQDc+e1eAhZ95g3qk\nXLeKEK9j+F+wuG/OrEAxBsscCxXRUB38QH6CFe3UxGoSMnBi+jEhicudo+ItpFOy\n7rPaYyRh4Pmu4QHcC83bNjp8NI6zTHrBmVuPqkxMn07GMAQav9ezBXj6umqTX4cU\ndsJUavJrJ3u7rT0lhBdiGrQ9zPbL07u2Kn+OXPAC3dKSf7G8TvwNAdry9esGSpi3\n8aa1myQYVZvAlsIBkbN3fb1wvDJE5czVhzwQ77V2t66jxeg0o9/2OZVH3CozD2Zj\nv28LHuW/jnQHtsQ0fUyQYRmHxNEVkW10GGM7fQwxzpxFFS1O/2XEnfMu7yBHZsgL\nSojfUct0FhLhEN/g/IINX9ZCVrzK5/De27CNjYE1cgYD/lTmQ0SyjfKVwwARAQAB\ntDFGZWRvcmEgKDMxKSA8ZmVkb3JhLTMxLXByaW1hcnlAZmVkb3JhcHJvamVjdC5v\ncmc+iQI+BBMBAgAoAhsPBgsJCAcDAgYVCAIJCgsEFgIDAQIeAQIXgAUCXGrkTQUJ\nEs8P/QAKCRBQyzkLPDNZxBmDD/90IFwAfFcQq5ENl7/o2CYQ9k2adTHbV5RoIOWC\n/o9I5/btn1y8WDhPOUNmsgbUqRqz6srlVplg+LkpIj67PVLDBwpVbCJC8o1fztd2\nMryVqdvu562WVhUorII+iW7nfqD0yv55nH9b/JR1qloUa8LpeKw84JgvxF5wVfyR\nid1WjI0DBk2taFR4xCfU5Tb262fbdFj5iB9xskP7oNeS29+SfDjlnybtlFoqr9UA\nnY1uvhBPkGmj45SJkpfP+L+kGYXVaUd29M/q/Pt46X1KDvr6Z0l8bSUEk3zfcNdj\nuEhtHBqSy1UPPAikGX1Q5wGdu7R7+mv/ARqfI6OC44ipoOMNK1Aiu6+slbPYphwX\nighSz9yYuG0EdWt7akfKR0R04Kuej4LXLWcxTR4l8XDzThYgPP0g+z0XQJrAkVhi\nSrzICeC3K1GPSiUtNAxSTL+qWWgwvQyAPNoPV/OYmY+wUxUnKCZpEWPkL79lh6CM\nbJx/zlrOMzRumSzaOnKW9AOliviH4Rj89OmDifBEsQ0CewdHN9ly6g4ZFJJGYXJ5\nHTb5jdButTC3tDfvH8Z7dtXKdC4iqJCIxj698Xn8UjVefZQ2nbv5eXcZLfHtvbNB\nTTv1vvBV4G7aiHKYRSj7HmxhLBZC8Y/nmFAemOoOYDpR5eUmPmSbFayoLfRsFXmC\nHLs7cw==\n=6hRW\n-----END PGP PUBLIC KEY BLOCK-----\n",
        }
    ],
    "fedora-32": [
        {
            "baseurl": "http://download.fedoraproject.org/pub/fedora/linux/releases/32/Everything/x86_64/os/",
            "gpgkey": "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmQINBFxq3QMBEADUhGfCfP1ijiggBuVbR/pBDSWMC3TWbfC8pt7fhZkYrilzfWUM\nfTsikPymSriScONXP6DNyZ5r7tgrIVdVrJvRIqIFRO0mufp9HyfWKDO//Ctyp7OQ\nzYw6NVthO/aWpyFfJpj6s4iZsYGqf9gByV8brBB8v8jEsCtVOj1BU3bMbLkMsRI9\n+WiLjDYyvopqNBQuIe8ogxSxpYdbUz6+jxzfvhRoBzWdjITd//Gjd90kkrBOMWkO\nLTqO133OD1WMT08G5NuQ4KhjYsVvSbBpfdkTcNuP8gBP9LxCQDc+e1eAhZ95g3qk\nXLeKEK9j+F+wuG/OrEAxBsscCxXRUB38QH6CFe3UxGoSMnBi+jEhicudo+ItpFOy\n7rPaYyRh4Pmu4QHcC83bNjp8NI6zTHrBmVuPqkxMn07GMAQav9ezBXj6umqTX4cU\ndsJUavJrJ3u7rT0lhBdiGrQ9zPbL07u2Kn+OXPAC3dKSf7G8TvwNAdry9esGSpi3\n8aa1myQYVZvAlsIBkbN3fb1wvDJE5czVhzwQ77V2t66jxeg0o9/2OZVH3CozD2Zj\nv28LHuW/jnQHtsQ0fUyQYRmHxNEVkW10GGM7fQwxzpxFFS1O/2XEnfMu7yBHZsgL\nSojfUct0FhLhEN/g/IINX9ZCVrzK5/De27CNjYE1cgYD/lTmQ0SyjfKVwwARAQAB\ntDFGZWRvcmEgKDMxKSA8ZmVkb3JhLTMxLXByaW1hcnlAZmVkb3JhcHJvamVjdC5v\ncmc+iQI+BBMBAgAoAhsPBgsJCAcDAgYVCAIJCgsEFgIDAQIeAQIXgAUCXGrkTQUJ\nEs8P/QAKCRBQyzkLPDNZxBmDD/90IFwAfFcQq5ENl7/o2CYQ9k2adTHbV5RoIOWC\n/o9I5/btn1y8WDhPOUNmsgbUqRqz6srlVplg+LkpIj67PVLDBwpVbCJC8o1fztd2\nMryVqdvu562WVhUorII+iW7nfqD0yv55nH9b/JR1qloUa8LpeKw84JgvxF5wVfyR\nid1WjI0DBk2taFR4xCfU5Tb262fbdFj5iB9xskP7oNeS29+SfDjlnybtlFoqr9UA\nnY1uvhBPkGmj45SJkpfP+L+kGYXVaUd29M/q/Pt46X1KDvr6Z0l8bSUEk3zfcNdj\nuEhtHBqSy1UPPAikGX1Q5wGdu7R7+mv/ARqfI6OC44ipoOMNK1Aiu6+slbPYphwX\nighSz9yYuG0EdWt7akfKR0R04Kuej4LXLWcxTR4l8XDzThYgPP0g+z0XQJrAkVhi\nSrzICeC3K1GPSiUtNAxSTL+qWWgwvQyAPNoPV/OYmY+wUxUnKCZpEWPkL79lh6CM\nbJx/zlrOMzRumSzaOnKW9AOliviH4Rj89OmDifBEsQ0CewdHN9ly6g4ZFJJGYXJ5\nHTb5jdButTC3tDfvH8Z7dtXKdC4iqJCIxj698Xn8UjVefZQ2nbv5eXcZLfHtvbNB\nTTv1vvBV4G7aiHKYRSj7HmxhLBZC8Y/nmFAemOoOYDpR5eUmPmSbFayoLfRsFXmC\nHLs7cw==\n=6hRW\n-----END PGP PUBLIC KEY BLOCK-----\n",
        }
    ],
    "rhel-8": [
        {"baseurl": "http://download.devel.redhat.com/released/RHEL-8/8.2.0/BaseOS/x86_64/os/"},
        {"baseurl": "http://download.devel.redhat.com/released/RHEL-8/8.2.0/AppStream/x86_64/os/"},
    ]
}


def compose_request(distro, koji):
    repositories = [repo for repo in DISTRO_REPOS[distro]]

    req = {
        "name": "name",
        "version": "version",
        "release": "release",
        "distribution": distro,
        "koji": {
            "server": koji,
            "task_id": 1
        },
        "image_requests": [{
            "architecture": "x86_64",
            "image_type": "qcow2",
            "repositories": repositories
        }]
    }

    return req


def main(distro):
    cr = compose_request(distro, "https://localhost:4343/kojihub")
    print(json.dumps(cr))

    r = requests.post("https://localhost/api/composer-koji/v1/compose", json=cr,
                      cert=("/etc/osbuild-composer/worker-crt.pem", "/etc/osbuild-composer/worker-key.pem"),
                      verify="/etc/osbuild-composer/ca-crt.pem")
    if r.status_code != 201:
        print("Failed to create compose")
        print(r.text)
        sys.exit(1)

    print(r.text)
    compose_id = r.json()["id"]

    while True:
        r = requests.get(f"https://localhost/api/composer-koji/v1/compose/{compose_id}",
                         cert=("/etc/osbuild-composer/worker-crt.pem", "/etc/osbuild-composer/worker-key.pem"),
                         verify="/etc/osbuild-composer/ca-crt.pem")
        if r.status_code != 200:
            print("Failed to get compose status")
            print(r.text)
            sys.exit(1)
        status = r.json()["status"]
        print(status)
        if status == "success":
            print("Compose worked!")
            print(r.text)
            break
        elif status == "failure":
            print("compose failed!")
            print(r.text)
            sys.exit(1)
        elif status != "pending" and status != "running":
            print(f"unexpected status: {status}")
            print(r.text)
            sys.exit(1)

        time.sleep(10)


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(f"usage: {sys.argv[0]} DISTRO", file=sys.stderr)
        sys.exit(1)
    main(sys.argv[1])
