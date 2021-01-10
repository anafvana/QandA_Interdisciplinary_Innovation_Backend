#!/bin/bash

ssh banana@ingrids.space 'killall server'
rsync server banana@ingrids.space:/home/banana/server
ssh banana@ingrids.space '/home/banana/server'

exit 0