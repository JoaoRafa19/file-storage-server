from concurrent.futures import ThreadPoolExecutor

import requests

PORTA = 8080
TIMEOUT = 2


def testar_ip(i):
    ip = f"192.168.1.{i}"
    url = f"http://{ip}:{PORTA}/"

    try:
        r = requests.get(url, timeout=TIMEOUT)

        if r.status_code == 200:
            print(f"[+] Encontrado: {url}")

            # Algumas interfaces de impressoras 3D têm palavras características
            conteudo = r.text.lower()
            if any(
                p in conteudo for p in ["octoprint", "mainsail", "fluidd", "klipper"]
            ):
                print(f"    Possível painel de impressão 3D detectado!")
    except:
        pass


with ThreadPoolExecutor(max_workers=20) as executor:
    executor.map(testar_ip, range(0, 100))
