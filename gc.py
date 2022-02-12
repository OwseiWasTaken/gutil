#! /usr/bin/python3.10
#imports
from util import *
from os.path import realpath

paths = {}

#main
def Main() -> int:
	global paths

	fl = argv[0]
	confile = '/'.join(fl.split('/')[:-1])+'/config.xmp'

	if not exists(confile):
		UseXmp(confile, {
			"paths":{},
	})

	paths = UseXmp(confile)["paths"]

	argvs = get("").list

	files = []
	config = {
		"run":False, # run, not compile
		"clear":False, # clear screen
		"kc":False, # keep cfile
		"lt":False, # light text
		"nb":False, # don't build cfile
		# other flags here
	}
	BuildArgs = []
	if get('-ba').exists:
		BuildArgs = get('-ba').list

	for arg in argvs:
		if len(arg):
			if arg[0] == '/' and arg[-1] == '/':
				config[arg[1:-1]] = True
			else:
				files.append(arg)

	if config["clear"]:
		ss("clear")

	for file in files:
		if ecode:=DoFileMain(file, config, BuildArgs):
			return ecode
	return 0

def DoFileMain(filename, config, BuildArgs) -> int:
	#filename = get(None).first
	run = config["run"]
	kc = config["kc"]
	bright = config["lt"]
	nb = config["nb"]
	# make text bright again
	if run and bright:
		stdout.write("\x1b[38;2;255;255;255m\n\x1b[1;1H")

	# usable input chek
	if not filename:
		fprintf(stderr, "input filename!\n")
		return 1
	elif not exists(filename):
		fprintf(stderr, "file \"{s}\" doesn't exist!\n", filename)
		return 2
	elif not isfile(filename):
		#TODO make dir buildable
		if "main.go" in ls(filename):
			cd(filename)
			filename="./main.go"
		else:
			fprintf(stderr, "can't find main.go in {s}\n", filename)
			return 3


	programname = '.'.join(filename.split('.')[:-1])

	with open(filename, 'r') as f:
		file, imports, includes, pack = CompFile(
			list(map(
				lambda x: x.replace('\n', '').strip(),
				f.readlines()
			)),
			{},
			True
		)

	newfile = MakeFile(file, imports, includes, pack)
	mvflname = filename

	# make "cfile"
	if '/' in filename:
		cname = filename.split('/')
		mvflname = cname[-1]
		cname = '/'.join(cname[:-1])+'/'+'c'+cname[-1]
	else:
		cname = 'c'+filename

	# write NewFile to cfile
	with open(cname, 'w') as f:
		f.writelines(map(lambda x: x+'\n', newfile))

	# build/run cfile
	if not get("--no-build").exists | nb:
		if run:
			if ss(f"go run {cname}"):
				fprintf(stderr, "could not run file {s}\n", filename)
		else:
			if ss(f"go build {' '.join(BuildArgs)} {cname}"):
				fprintf(stderr, "could not compile file {s}\n", filename)
				if not kc:
					ss(f"rm {cname}")
				return 1
			ss(f"mv c{mvflname[:-3]} {mvflname[:-3]}")

		if not kc:
			ss(f"rm {cname}")
	return 0

def CompFile(file: list[str], includes={} , RetPack = True) -> tuple[list[str], list[str], list[str], Optional[str]]:
	global paths

	FILE = []
	imporing = False
	imports = set([])
	pack = ""
	for line in file:
		if line[:7] == "package":
			pack = line
		elif line == ")" and imporing:
			imporing = False
		elif imporing:
			imports.add(line)
		elif line[:6] == "import" and line[7] == '(':
			imporing = True
		elif line[:7] == "include":
			includename = line[9:][:-1]
			if includename in paths.keys():
				includename = paths[includename]
			else:
				if not includename.endswith('.go'):
					includename+=".go"

				if not exists(includename):
					if exists("../"+includename):
						includename = "../"+includename
					else:
						if includename == "gutil.go":
							if exists("gutil/"+includename):
								includename = "gutil/"+includename
						else:
							fprintf(stderr, "can't find included file {s}\n", includename)
			if exists(includename):
				if (includename in includes.keys()):continue
				#FL = []
				with open(includename, 'r') as f:
					FL, _imports, _includes = CompFile(
						list(map(
							lambda x: x.replace('\n', '').strip(),
							f.readlines()
						)),includes , False
					)
					imports = set([*imports, *_imports])
					includes[includename] = FL
					includes = includes | _includes
				#FILE+=FL[1:]
			else:
				fprintf(stderr, "No Such File \"{s}\"\n", includename)
		else:
			FILE.append(line)
	if RetPack:
		return FILE, imports, includes, pack
	else:
		return FILE, imports, includes

def MakeFile(file, imports, includes, pack):
	FL = []
	for n, fl in includes.items():
		FL.append("// include %s"%n)
		FL = FL+fl
	file = FL+file
	file=[pack, "import ("]+[x for x in list(imports)]+[")"]+file
	return file

#start
if __name__ == '__main__':
	start = tm()
	ExitCode = Main()

	if get('--debug').exists:
		if not ExitCode:
			printl("%scode successfully exited in " % COLOR.green)
		else:
			printl("%scode exited with error %d in " % (COLOR.red,ExitCode))
		print("%.3f seconds%s" % (tm()-start,COLOR.nc))
	exit(ExitCode)
