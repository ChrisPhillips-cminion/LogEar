FROM debian
COPY ./bugle*  /bugle
ENV namespace default
ENV podname podname
ENV username user
ENV password pass
ENTRYPOINT /bugle -podname $podname -namespace $namespace -username $username -password $password
