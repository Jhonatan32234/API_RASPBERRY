import cv2
import time
import requests
from datetime import datetime

ZONA = "zona_1"

def detectar_personas():
    cap = cv2.VideoCapture(0)
    if not cap.isOpened():
        print("❌ No se pudo abrir la cámara.")
        return

    conteo_personas = 0

    while True:
        ret, frame = cap.read()
        if not ret:
            break

        # Aquí pones la lógica de detección (por ejemplo, MobileNetSSD o fondo)
        # Por simplicidad, simulamos detección
        # TODO: reemplaza con tu modelo real

        # Simulación: detecta 1 persona cada 5 segundos
        conteo_personas += 1
        print(f"👤 Persona detectada. Total: {conteo_personas}")

        time.sleep(5)

        # Puedes definir alguna condición para salir y enviar datos
        if conteo_personas >= 10:
            break

    cap.release()

    # Preparar datos para enviar
    now = datetime.now()
    payload = [{
        "visitantes": conteo_personas,
        "hora": now.strftime("%H:%M:%S"),
        "fecha": now.strftime("%Y-%m-%d"),
        "zona": ZONA,
        "enviado": False
    }]

    try:
        print(f"📤 Enviando datos a API Go: {payload}")
        r = requests.post("http://localhost:8080/visitas", json=payload)
        print(f"✅ Enviado con código {r.status_code}")
    except Exception as e:
        print(f"❌ Error al enviar datos: {e}")

if __name__ == "__main__":
    detectar_personas()
