FROM debian
COPY ./LogEar*  /LogEar
ENV namespace default
ENV podname podname
ENV username user
ENV password pass
ENTRYPOINT /LogEar -podname $podname -namespace $namespace -username $username -password $password
