import cv2
import time
from datetime import datetime

INTERVALO_SEGUNDOS = 5

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
persona_en_cuadro = False

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

        conteo_personas_frame = 0

        for i in range(detections.shape[2]):
            confidence = detections[0, 0, i, 2]
            if confidence > 0.5:
                idx = int(detections[0, 0, i, 1])
                if CLASSES[idx] == "person":
                    conteo_personas_frame += 1
                    box = detections[0, 0, i, 3:7] * [w, h, w, h]
                    (startX, startY, endX, endY) = box.astype("int")
                    cv2.rectangle(frame, (startX, startY), (endX, endY), (0, 255, 0), 2)

        if conteo_personas_frame > 0 and not persona_en_cuadro:
            nuevas_personas += 1
            persona_en_cuadro = True
            print("🟢 Una persona entró al cuadro.")

        if conteo_personas_frame == 0 and persona_en_cuadro:
            persona_en_cuadro = False

        cv2.putText(frame, f'Personas en cuadro: {conteo_personas_frame}', (10, 30),
                    cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), 2)
        cv2.putText(frame, f'Total entradas: {nuevas_personas}', (10, 60),
                    cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), 2)

        cv2.imshow("Detección", frame)

        if time.time() - inicio_intervalo >= INTERVALO_SEGUNDOS:
            timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            resultado = f"{timestamp} - Personas entraron: {nuevas_personas}\n"
            print(f"🕒 Guardando en TXT: {resultado.strip()}")

            with open('registro_personas.txt', 'a', encoding='utf-8') as archivo:
                archivo.write(resultado)

            nuevas_personas = 0
            inicio_intervalo = time.time()

        if cv2.waitKey(1) & 0xFF == ord('q'):
            print("Saliendo por tecla 'q'...")
            break

finally:
    cap.release()
    cv2.destroyAllWindows()


""" import cv2
import time
from datetime import datetime

INTERVALO_SEGUNDOS = 10  # Usa 3600 para 1 hora

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
personas_intervalo = 0

archivo = open('conteo_intervalos.txt', 'a', encoding='utf-8')

try:
    while True:
        ret, frame = cap.read()
        if not ret:
            print("No se pudo leer frame de la cámara.")
            break

        h, w = frame.shape[:2]
        blob = cv2.dnn.blobFromImage(cv2.resize(frame, (300, 300)),
                                     0.007843, (300, 300), 127.5)

        net.setInput(blob)
        detections = net.forward()

        conteo_personas_frame = 0

        for i in range(detections.shape[2]):
            confidence = detections[0, 0, i, 2]
            if confidence > 0.5:
                idx = int(detections[0, 0, i, 1])
                if CLASSES[idx] == "person":
                    conteo_personas_frame += 1
                    box = detections[0, 0, i, 3:7] * [w, h, w, h]
                    (startX, startY, endX, endY) = box.astype("int")
                    cv2.rectangle(frame, (startX, startY), (endX, endY), (0, 255, 0), 2)

        if conteo_personas_frame > personas_intervalo:
            personas_intervalo = conteo_personas_frame

        cv2.putText(frame, f'Personas detectadas: {conteo_personas_frame}', (10, 30),
                    cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), 2)
        cv2.imshow("Detección", frame)

        # Verificar si terminó el intervalo
        if time.time() - inicio_intervalo >= INTERVALO_SEGUNDOS:
            timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            resultado = f"{timestamp} - Personas detectadas: {personas_intervalo}\n"
            print(f"🕒 Guardando en TXT: {resultado.strip()}")

            archivo.write(resultado)
            archivo.flush()

            inicio_intervalo = time.time()
            personas_intervalo = 0

        # Salir con 'q'
        if cv2.waitKey(10) & 0xFF == ord('q'):
            print("Saliendo por tecla 'q'...")
            break
finally:
    archivo.close()
    cap.release()
    cv2.destroyAllWindows()
 """