# set base image (host OS)
FROM python:3.8

# set the working directory in the container
RUN mkdir -p /opt
WORKDIR /opt

# copy the content of the local src directory to the working directory
COPY . .

# install dependencies
RUN pip install -r requirements.txt

CMD ["python3", "main.py"]

