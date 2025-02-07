import subprocess
import time

from watchdog.events import FileSystemEventHandler
from watchdog.observers import Observer


class ReloadHandler(FileSystemEventHandler):
    """
    ファイルの変更を監視し、gRPCサーバーを再起動するハンドラー
    """

    def __init__(self, command, watch_extensions=None):
        self.command = command
        self.process = None
        self.last_restart_time = 0
        self.restart_interval = 1  # 再起動間隔（秒）
        self.watch_extensions = watch_extensions or [".py"]  # 監視対象の拡張子
        self.start_server()

    def start_server(self):
        """
        サーバーを起動する
        """
        if self.process:
            print("Stopping existing server...")
            self.process.terminate()
            self.process.wait()  # 完全終了を待つ
        print("Starting server...")
        self.process = subprocess.Popen(self.command, shell=True)

    def on_modified(self, event):
        """
        ファイルが変更されたときにサーバーを再起動する
        """
        if event.is_directory:
            return  # ディレクトリ変更は無視
        if not any(event.src_path.endswith(ext) for ext in self.watch_extensions):
            return  # 対象外のファイルは無視

        # 再起動間隔を制御
        current_time = time.time()
        if current_time - self.last_restart_time < self.restart_interval:
            return
        self.last_restart_time = current_time

        print(f"File changed: {event.src_path}. Reloading server...")
        self.start_server()


if __name__ == "__main__":
    grpc_command = "python src/server/main.py"

    # 監視対象の拡張子
    watch_extensions = [".py"]

    # ハンドラーとオブザーバーを設定
    event_handler = ReloadHandler(grpc_command, watch_extensions)
    observer = Observer()
    observer.schedule(event_handler, path="src", recursive=True)
    observer.start()

    try:
        print("Watching for file changes in 'src'...")
        observer.join()
    except KeyboardInterrupt:
        print("Stopping observer...")
        observer.stop()
    observer.join()
