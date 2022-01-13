#! /usr/bin/python3.10
#imports
from util import *
# TODO
# store already compiled files in /tmp and hash them
# compare compiling files to the already compiled ones
# how? idk
# why? speed prolly

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
		'r':False,#run
		# other flags here
	}
	BuildArgs = []
	if get('-ba').exists:
		BuildArgs = get('-ba').list

	for arg in argvs:
		if arg[0] == '/' and arg[-1] == '/':
			config[arg[1]] = True
		else:
			files.append(arg)

	for file in files:
		if ecode:=DoFileMain(file, config, BuildArgs):
			return ecode
	return 0

def DoFileMain(filename, config, BuildArgs) -> int:
	#filename = get(None).first
	run = config['r']

	# usable input chek
	if not filename:
		fprintf(stderr, "input filename!\n")
		return 1
	elif not exists(filename):
		fprintf(stderr, "file \"{s}\" doesn't exist!\n", filename)
		return 2
	elif not isfile(filename):
		#TODO make dir buildable
		fprintf(stderr, "\"{s}\" is a directory, now a file\n", filename)
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
	if '/' in filename:
		cname = filename.split('/')
		mvflname = cname[-1]
		cname = '/'.join(cname[:-1])+'/'+'c'+cname[-1]
	else:
		cname = 'c'+filename
	with open(cname, 'w') as f:
		f.writelines(map(lambda x: x+'\n', newfile))
	if not get("--no-build").exists:
		if run:
			if ss(f"go run {cname}"):
				fprintf(eout, "could not run file {s}\n", filename)
		else:
			if ss(f"go build {' '.join(BuildArgs)} {cname}"):
				fprintf(eout, "could not compile file {s}\n", filename)
				if not get('-ke').exists:
					ss(f"rm {cname}")
				return 1
			ss(f"mv c{mvflname[:-3]} {mvflname[:-3]}")
		if not get('-ke').exists:
			ss(f"rm {cname}")
	return 0

def CompFile(file: list[str], includes={} , RetPack = True) -> list[str]:
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
