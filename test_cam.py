import cv2

cap = cv2.VideoCapture(0)
if not cap.isOpened():
    print("No se pudo abrir la cámara 0")
    exit()

cv2.namedWindow("Detección", cv2.WINDOW_NORMAL)

while True:
    ret, frame = cap.read()
    if not ret:
        print("No se pudo leer frame")
        break

    cv2.imshow("Detección", frame)

    if cv2.waitKey(10) & 0xFF == ord('q'):
        break

cap.release()
cv2.destroyAllWindows()

