int tambah(int a, int b) {
    return a + b;
}

int kurang(int a, int b) {
    return a - b;
}

int kali(int a, int b) {
    return a * b;
}

int bagi(int a, int b, int *rem) {
    int ret = a / b;
    if(rem) {
        *rem = a % b;
    }
    return ret;
}
