# coding:utf8

import os
import os.path
import re
import shutil



def write_file(filename, txt):
    with open(filename, 'w') as f:
        f.write(txt)
        f.close()


def read_file(filename):
    with open(filename, 'r') as content_file:
        content = content_file.read()
        content_file.close()
        return content

def get_protocol_files(rootDir,fileregs,isfind):
	files = []
	for lists in os.listdir(rootDir):
		path = os.path.join(rootDir, lists)
		if os.path.isdir(path):
			dirpath = path+"/protocol"
			for path_ in os.listdir(dirpath):
				fpath = os.path.join(dirpath, path_)
				
				if os.path.isdir(fpath):
					continue
				
				find = False
				for reg in fileregs:
					if fpath.find(reg) != -1:
						find = True
						break
					
				if  find and isfind:
					files.append(fpath)
				elif not find and not isfind:
					files.append(fpath)
				
	return files

# 合成CMD
def merge_cmd():
    file_lists = [
	"../../src/server_gateway/protocol/gateway_cmd.proto",
                "../../src/server_login/protocol/login_cmd.proto",
				  "../../src/server_dba/protocol/dba_cmd.proto",
    #             "../../src/game_protocol/game_cmd.proto",
	#		      "../../src/server_game_hall/protocol/game_hall_cmd.proto",
                  "../../src/server_console/protocol/console_cmd.proto",
	#			  "../../src/match/protocol/game_match_cmd.proto",
				]
 #   protofiles = get_protocol_files("../../src/server_games/",["cmd.proto"],True)
 #   for filename in protofiles:
 #       file_lists.append(filename)
    buffer = "package protocol;\n\nenum CMD {"
    for files in file_lists:
        # print (files)
        temp = read_file(files)
        pattern = re.compile("package\s*protocol\s*;")
        temp = pattern.sub('', temp, 1)
        pattern = re.compile("enum\s*CMD\s*\{")
        temp = pattern.sub('', temp, 1)
        pattern = re.compile("\}")
        temp = pattern.sub('', temp, 1)
        buffer += temp

    buffer += "}"
    write_file("../../protocol/cmd.proto", buffer)


def merge_status_code():
    buffer = "package protocol;\n\nenum StatusCode{\n"
    code_lists = [
	"../../src/server_login/protocol/login_code.proto",
				  "../../src/game_protocol/game_code.proto",
	#			  "../../src/match/protocol/game_match_code.proto",
				]
 #   protofiles = get_protocol_files("../../src/server_games/",["code.proto"],True)
 #   for filename in protofiles:
 #       code_lists.append(filename)
    for files in code_lists:
        # print (files)
        temp = read_file(files)

        pattern = re.compile("package\s*protocol\s*;")
        temp = pattern.sub('', temp, 1)

        pattern = re.compile("enum\s*StatusCode\s*\{")
        temp = pattern.sub('', temp, 1)

        pattern = re.compile("\}")
        temp = pattern.sub('', temp, 1)
        # 删除空行
        # pattern = re.compile(r'\r.\n')
        # pattern = re.compile(r'^\s*$')
        # temp = pattern.sub('', temp)

        buffer += temp

    buffer += "}"
    write_file("../../protocol/statuscode.proto", buffer)


def copy_other_file():
    # print ("copy other files")
    dirs = [
	"../../src/server_gateway/protocol/",
            "../../src/server_login/protocol/",
			"../../src/server_dba/protocol/",
     #       "../../src/game_protocol/",
     #       "../../src/match/protocol/",
	#		"../../src/server_game_hall/protocol/",
            "../../src/server_console/protocol/",
			]
    for folder in dirs:
        for lists in os.listdir(folder):
            path_ = os.path.join(folder, lists)
            if os.path.isdir(path_):
                print ("is dir",path_)
            else:
                if path_.find("cmd.proto") !=-1 or path_.find("code.proto")!=-1:
                    continue
                    # print ("don't copy",path_)
                else:
                    # print path_
                    path_dir,path_name = os.path.split(path_)
                    # os.path.dirname(existGDBPath)
                    shutil.copyfile(path_, "../../protocol/" + path_name)

  #  protofiles = get_protocol_files("../../src/server_games/",["code.proto","cmd.proto"],False)
  #  for filename in protofiles:
  #      path_dir,path_name = os.path.split(filename)
  #      shutil.copyfile(filename, "../../protocol/" + path_name)

def copy_config():
    shutil.copyfile("../../src/lib/public_config/public_config.toml", "../../bin/public_config.toml")


if __name__ == '__main__':
    copy_config()
    print ("merge command")
    merge_cmd()
    print ("merge status code")
    merge_status_code()
    print ("copy file.")
    copy_other_file()
    print (u"merge finished ok")
    # print (os.system("./bat/generate.sh"))
    # 遍历目录
    # 合成cmd.proto
    # allfile = dirlist("./",[])
    # doPro("cmd.proto",allfile)
