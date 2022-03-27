#! /usr/bin/python3.10
#imports
from util import *
from os.path import realpath, abspath

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
	if get("-l").exists:
		ll = get("-l").list
		np = {}
		for l in ll:
			if ':' in l:
				np[l.split(':')[0]] = l.split(':')[1]
			else:
				fprintf(stderr, "-l error, can't find ':' division on `{s}`", l)
				exit(1)
		paths = paths | np

	argvs = get("").list

	files = []
	config = {
		"run":False, # run, not compile
		"clear":False, # clear screen
		"kc":False, # keep cfile
		"lt":False, # light text
		"nb":False, # don't build cfile
		"bnr":False, # build and run
		# other flags here
	}
	BuildArgs = []
	RunArgs = []
	if get('--ba').exists:
		BuildArgs = get('--ba').list
	if get('--ra').exists:
		RunArgs = get('--ra').list

	for arg in argvs:
		if len(arg):
			if arg[0] == '/' and arg[-1] == '/':
				config[arg[1:-1]] = True
			else:
				if '.' in arg:
					files.append(arg)
				else:
					if exists(arg+".go") and not exists(arg):
						files.append(arg+".go")
	if len(files) == 0:
		files = ["."]

	if config["clear"]:
		cmd("clear")

	for file in files:
		if ecode:=DoFileMain(file, config, BuildArgs, RunArgs):
			return ecode
	return 0

def DoFileMain(filename, config, BuildArgs, RunArgs) -> int:
	#filename = get(None).first
	run = config["run"]
	kc = config["kc"]
	bright = config["lt"]
	nb = config["nb"]
	bnr = config["bnr"]
	# make text bright again
	if run and bright:
		stdout.write("\x1b[38;2;255;255;255m\n\x1b[1;1H")

	# usable input chek
	if not exists(filename):
		fprintf(stderr, "file \"{s}\" doesn't exist!\n", filename)
		return 2
	elif not isfile(filename):
		if "main.go" in ls(filename):
			cd(filename)
			filename="./main.go"
		else:
			fprintf(stderr, "can't find main.go in {s}\n", filename)
			return 3

	now = pwd()
	env = false
	if "/" in filename:
		env = true
		cd( '/'.join(filename.split("/")[:-1])+"/" )
		filename = filename.split("/")[-1]


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

	if env:
		cd(now) # old wd

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
	if not nb:
		if run:
			ra = ' - '+''.join(RunArgs)+' '
			if cmd(f"go run {cname} {ra}"):
				fprintf(stderr, "could not run file {s}\n", filename)
		else:
			ba = ' '+''.join(DoAll(lambda x: " -"+x, BuildArgs))+' '
			if _:=cmd(f"go build {cname}{ba}"):
				fprintf(stderr, "could not compile file {s}\n", filename)
				if not kc:
					cmd(f"rm {cname}")
				return 1
			cmd(f"mv c{mvflname[:-3]} {mvflname[:-3]}") # renames compiled c{file} -> {file}
			if not _ and bnr: # build and run
				cmd(f"./{mvflname[:-3]}")
		if not kc:
			cmd(f"rm {cname}")
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
		elif line[0:2] == "#!":
			continue
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
				if not includename.endswith('.go') and exists(includename+".go"):
					includename+=".go"
				if not exists(includename):
					if exists("../"+includename):
						includename = "../"+includename
					else:
						fprintf(stderr, "can't find included file {s}\n", includename)
						exit(1)
			if exists(includename):
				if (abspath(includename) in includes.keys()):continue
				#FL = []
				if not isfile(includename):
					if exists(includename+"/main.go"):
						# remember .
						p = pwd()
						# goto included folder
						cd(includename)

						# include folder/main.go
						with open(includename+"/main.go", 'r') as f:
							FL, _imports, _includes = CompFile(
								list(map(
									lambda x: x.replace('\n', '').strip(),
									f.readlines()
								)),includes , False
							)
							FL.insert(0, "// include %s/" % includename)
							imports = set([*imports, *_imports])
							includes[abspath(includename)] = FL
							includes = includes | _includes

						# comeback
						cd(p)
					else:
						fprintf(
							stderr,
							"can't include dir {s}, can't find {s}/main.go\n",
							includename, includename
						)
						exit(1)
				else:
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
				exit(1)
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
