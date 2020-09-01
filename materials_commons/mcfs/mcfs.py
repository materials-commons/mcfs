from fs import ResourceType
from fs.base import FS
from fs.info import Info
from fs.mode import Mode
from fs.subfs import SubFS
from fs.time import datetime_to_epoch
from materials_commons.api import Client, MCAPIError


class MCFS(FS):
    def __init__(self, project_id, apitoken, base_url):
        super(MCFS, self).__init__()
        self._project_id = project_id
        self._c = Client(apitoken, base_url)

    def listdir(self, path):
        self.check()
        dir_listing = self._c.list_directory_by_path(self._project_id, path)
        return [d.path for d in dir_listing]

    def makedir(self, path, permissions=None, recreate=False):
        self.check()
        self._c.create_directory_by_path(self._project_id, path)
        return SubFS(self, path)

    def openbin(self, path, mode="r", buffering=-1, **options):
        self.check()
        _mode = Mode(mode)
        _mode.validate_bin()
        if _mode.create:
            def on_close_create(mcfile):
                try:
                    mcfile.raw.seek(0)
                    self._c.upload_file(self._project_id, 1, )
                finally:
                    pass

    def remove(self, path):
        pass

    def removedir(self, path):
        pass

    def setinfo(self, path, info):
        pass

    def isdir(self, path):
        try:
            return self.getinfo(path).is_dir
        except MCAPIError:
            return False

    def getinfo(self, path, namespaces=None):
        if path == "/":
            return Info({
                "basic": {"name": "", "is_dir": True},
                "details": {"type": int(ResourceType.directory)}
            })
        f = self._c.get_file_by_path(self._project_id, path)
        is_dir = f.mime_type == "directory"
        return Info({
            "basic": {"name": f.name, "is_dir": is_dir},
            "modified": datetime_to_epoch(f.mtime),
            "size": f.size,
            "type": int(ResourceType.directory if is_dir else ResourceType.file)
        })
