
import cv2
import time
from datetime import datetime
import numpy as np

INTERVALO_SEGUNDOS = 5
UMBRAL_DISTANCIA = 50
TIEMPO_CENTRO_SEGUIMIENTO = 3  # segundos

net = cv2.dnn.readNetFromCaffe(
    'MobileNetSSD_deploy.prototxt',
    'MobileNetSSD_deploy.caffemodel'
)

CLASSES = ["background", "aeroplane", "bicycle", "bird", "boat",
           "bottle", "bus", "car", "cat", "chair", "cow", "diningtable",
           "dog", "horse", "motorbike", "person", "pottedplant",
           "sheep", "sofa", "train", "tvmonitor"]

cap = cv2.VideoCapture(0)
if not cap.isOpened():
    print("No se pudo abrir la cámara.")
    exit()

print("Iniciando detección...")

cv2.namedWindow("Detección", cv2.WINDOW_NORMAL)

inicio_intervalo = time.time()
nuevas_personas = 0
centros_recientes = []  # [(x, y, timestamp)]

try:
    while True:
        ret, frame = cap.read()
        if not ret:
            print("No se pudo leer frame.")
            break

        h, w = frame.shape[:2]

        blob = cv2.dnn.blobFromImage(cv2.resize(frame, (300, 300)),
                                     0.007843, (300, 300), 127.5)
        net.setInput(blob)
        detections = net.forward()

        tiempo_actual = time.time()
        nuevos_centros = []

        for i in range(detections.shape[2]):
            confidence = detections[0, 0, i, 2]
            if confidence > 0.5:
                idx = int(detections[0, 0, i, 1])
                if CLASSES[idx] == "person":
                    box = detections[0, 0, i, 3:7] * np.array([w, h, w, h])
                    (startX, startY, endX, endY) = box.astype("int")

                    centroX = int((startX + endX) / 2)
                    centroY = int((startY + endY) / 2)
                    nuevo_centro = (centroX, centroY)

                    # Dibujar la caja y centro
                    cv2.rectangle(frame, (startX, startY), (endX, endY), (0, 255, 0), 2)
                    cv2.circle(frame, nuevo_centro, 5, (255, 0, 0), -1)

                    # Verificar si ya fue contado recientemente
                    ya_contado = False
                    for cx, cy, t in centros_recientes:
                        dist = np.linalg.norm(np.array((cx, cy)) - np.array(nuevo_centro))
                        if dist < UMBRAL_DISTANCIA:
                            ya_contado = True
                            break

                    if not ya_contado:
                        nuevas_personas += 1
                        print("🟢 Nueva persona detectada")
                        centros_recientes.append((centroX, centroY, tiempo_actual))

                    nuevos_centros.append((centroX, centroY, tiempo_actual))

        # Limpiar centros viejos
        centros_recientes = [(x, y, t) for (x, y, t) in centros_recientes if tiempo_actual - t < TIEMPO_CENTRO_SEGUIMIENTO]

        cv2.putText(frame, f'Total entradas: {nuevas_personas}', (10, 30),
                    cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), 2)

        cv2.imshow("Detección", frame)

        # Guardar en archivo si se cumple el intervalo
        if tiempo_actual - inicio_intervalo >= INTERVALO_SEGUNDOS:
            timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            resultado = f"{nuevas_personas}\n"
            print(f"🕒 Guardando en TXT: {resultado.strip()}")

            with open('registro_personas.txt', 'a', encoding='utf-8') as archivo:
                archivo.write(resultado)

            nuevas_personas = 0
            inicio_intervalo = tiempo_actual

        if cv2.waitKey(1) & 0xFF == ord('q'):
            print("Saliendo por tecla 'q'...")
            break

finally:
    cap.release()
    cv2.destroyAllWindows()
