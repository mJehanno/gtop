- /proc/[nombre]/fd => une entrée pour chaque fichier que le processus a ouvert. Chaque entrée a le descripteur du fichier pour nom, et est représentée par un lien symbolique sur le vrai fichier. Ainsi, 0 correspond à l'entrée standard, 1 à la sortie standard, 2 à la sortie d'erreur
- /proc/[nombre]/stat => Informations sur l'état du processus. Ceci est utilisé par ps(1). La définition se trouve dans /usr/src/linux/fs/proc/array.c
- /proc/[nombre]/status => format human readable de stat et statm (notamment mémoire)
- /proc/loadavg => Les trois premiers champs de ce fichier sont les charges système indiquant le nombre de tâches dans la file d'attente d'exécution (état R) ou en attente d'entrées-sorties disque (état D) moyennées sur 1, 5, et 15 minutes
- /proc/net => informations réseaux (maybe check /proc/net/dev ?)
- /etc/lsd_release => information OS - OK
- /proc/version => version noyau linux - OK


Use [this repo](https://github.com/knipferrc/teacup#image) to display images of distribution in basic information. It needs either to check connectivity before retrieving assets on official website (and maybe resize) or embed images as assez (maybe using `go:embed`)



Revoir le modèle de données : 

Stats - For Linux/ For Mac / For Windows (only need to change suffix in filename, no need to create an interface and strcut for each OS)
	|
	- Metrics
	- SystemInformation


Basic Information Tab :

- Information about OS (distrib + kernel version + uptime) - ok
- Information about current user (uid, groups) - ok
- Information about ram (total and free) - ok
- Information about swap (total and free) - ok
- Information about CPU (model, and nomber of core)
- Information about GPU Model (if possible)
- Information about network (current public IP Address or error like "not connected to network")
- Information about disk size (total and free)

Usage Tab :

- Ram and swap usage (graph view) - ok
- cpu usage (graph view) - need to fix porgress bar
- disk usage (graph view)

Process Manager Tab :

- List of process (pid, name, path, mem-usage,cpu-usage) - miss cpu usage
- Tree view of process and sub-process
- Send signal to process (ex : SIGTERM)
- Filter process on user, name, pid - ok
- Sort process by mem-usage or cpu-usage or name (?) - miss cpu usage
- Port used by process

Network Manager Tab : 

- List all network card with IP
- Data received / Data Sent
