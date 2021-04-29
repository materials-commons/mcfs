from pathlib import Path


class MCProjectFileCache(object):
    def __init__(self, project):
        self.project = project
        self.cache_dir = str(Path.home() / '.mcfs' / project.uuid)
        Path(self.cache_dir).mkdir(exist_ok=True)
