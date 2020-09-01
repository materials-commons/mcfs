
from fs.opener import Opener


class MCFSOpener(Opener):
    protocols = ["mc"]

    def open_fs(self, fs_url, parse_result, writeable, create, cwd):
        pass