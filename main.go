package main

import (
	"image/png"
	"image"
	"os"
	"strconv"
	"fmt"
	"github.com/EyciaZhou/geo.go"
	"math"
	"image/color"
	"github.com/everdev/mack"
	"time"
	"encoding/json"
	"errors"
)

func getImgByFilename(fn string) (image.Image, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("can't open file : %s , reason : %s\n", fn, err.Error())
	}
	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("error when decode the png file %s, reason : %s\n", fn, err.Error())
	}
	return img, nil
}

func writeImg(fn string, img image.Image) error {
	f, err := os.Create(fn)
	if err != nil {
		return fmt.Errorf("can't open file : %s , reason : %s\n", fn, err.Error())
	}
	err = png.Encode(f, img)
	if err != nil {
		return fmt.Errorf("error when encode the png file, reason : %s\n", err.Error())
	}
	return nil
}

func mergePicturesAndRotate(path, pfix, topath, topfix string, l, r, t, b int, ang float64, back color.Color) error {
	//if len(os.Args) < 6 {
	//	fmt.Printf("Usage: ./mergepic pfix l r t b [-ro ang [-back co]] [-sel x0,y0,x1,y1] \n")
	//	return
	//}

	//etcArgs := map[string]string{}

	//pfix := os.Args[1]
	//l := getIntAndFail(os.Args[2])
	//r := getIntAndFail(os.Args[3])
	//t := getIntAndFail(os.Args[4])
	//b := getIntAndFail(os.Args[5])

	//for i := 6; i < (len(os.Args)-1)/2*2+1; i++ {
	//	etcArgs[os.Args[i]] = os.Args[i+1]
	//}

	cav := image.NewRGBA(image.Rect(0, 0, (r-l+1)*550, (b-t+1)*550))

	for i := l; i <= r; i++ {
		for j := t; j <= b; j++ {
			img, err := getImgByFilename(fmt.Sprintf(path+"%s_%d_%d.png", pfix, i, j))

			if err != nil {
				return err
			}

			for ii := 0; ii < 550; ii++ {
				for jj := 0; jj < 550; jj++ {
					cav.Set((i-l)*550+ii, (j-t)*550+jj, img.At(ii, jj))
				}
			}
		}
	}

	if ang != 0 {
		mat := geo.Rotate(ang/180*math.Pi)
		xs := [4]float64{0, 0, float64((r-l+1)*550), float64((r-l+1)*550)}
		ys := [4]float64{float64((b-t+1)*550), 0, float64((b-t+1)*550), 0}
		//rotate with center and clockwise
		for i := 0; i < 4; i++ {
			xs[i], ys[i] = mat.Apply(xs[i], ys[i])
		}
		mxx := math.Max(math.Max(xs[0], xs[1]), math.Max(xs[2], xs[3]))
		mxy := math.Max(math.Max(ys[0], ys[1]), math.Max(ys[2], ys[3]))
		mix := math.Min(math.Min(xs[0], xs[1]), math.Min(xs[2], xs[3]))
		miy := math.Min(math.Min(ys[0], ys[1]), math.Min(ys[2], ys[3]))

		fmt.Printf("%lf %lf %lf %lf\n", mxx, mxy, mix, miy)

		rect := image.Rect(0, 0, int(mxx)-int(mix), int(mxy)-int(miy))

		tmpimg := cav
		cav = image.NewRGBA(rect)
		//back := color.RGBA{0, 0, 0, 0xff}
		//if bg, ok := etcArgs["-back"]; ok {
		//	back = getRGBAndFail(bg)
		//}
		//fmt.Printf("%lf %lf\n", mix, miy)
		//do rotate
		mat = mat.Move(-mix, -miy).Inv()

		//fmt.Printf("%v\n", mat.Inv())

		dx := float64(tmpimg.Bounds().Dx())
		dy := float64(tmpimg.Bounds().Dy())
		for i := 0; i < cav.Bounds().Dx(); i++ {
			for j := 0; j < cav.Bounds().Dy(); j++ {
				x1, y1 := mat.Apply(float64(i), float64(j))
				x1 = x1 / dx
				y1 = y1 / dy
				cav.Set(i, j, GetPix(tmpimg, x1, y1, back))
			}
		}
	}

	return writeImg(topath+topfix+".png", cav)
}

func setDesktopPicture(fn string) error {
	return mack.Tell("Finder", fmt.Sprintf("set desktop picture to POSIX file \"%s\"", fn))
}

func downloadFile(url, fn string) error {
	f, e := os.Create(fn)
	if e != nil {
		return e
	}
	bs, e := Get(url)
	if e != nil {
		return e
	}
	_, e = f.Write(bs)
	return e
}

func getNewestDate() (date string, e error){
	//return "2015-01-02 14:15:16", nil

	bs, e := Get("http://himawari8-dl.nict.go.jp/himawari8/img/D531106/latest.json?uid=" + strconv.FormatInt(time.Now().Unix(), 10))

	if e != nil {
		return "", e
	}

	var m map[string]string
	e = json.Unmarshal(bs, &m)
	if e != nil {
		return "", e
	}
	if _, ok := m["date"]; !ok {
		return "", errors.New("no date in json")
	}
	return m["date"], nil
}

func downloadPictures(left, right, top, bottom int, filenamePrefix string, d, data string) error {
	a := data

	for i := left; i <= right; i++ {
		for j := top; j <= bottom; j++ {
			time.Sleep(1*time.Second)
			url := fmt.Sprintf("http://himawari8-dl.nict.go.jp/himawari8/img/D531106/%s/550/%s/%s/%s/%s%s%s_%d_%d.png", d, a[0:4], a[5:7], a[8:10], a[11:13], a[14:16], a[17:19], i, j)
			fmt.Println(url)
			e := downloadFile(url, fmt.Sprintf(filenamePrefix+"_%d_%d.png", i, j))
			if e != nil {
				return e
			}
		}
	}
	return nil
}

var (
	lst_da string
)

var (
	d string = "4d"
)

func downloadOnce(date string, setWallpaper bool) error {
	wd, e := os.Getwd()
	if e != nil {
		return e
	}

	pathl := wd+"/ep/tmp/"
	fnl := fmt.Sprintf("%s_550_%s_%s_%s_%s%s%s", d, date[0:4], date[5:7], date[8:10], date[11:13], date[14:16], date[17:19])
	pathr := wd+"/ep/"
	fnr := fnl

	if _, err := os.Stat(pathr+fnr+".png"); err == nil {
		return errors.New("File exists")
	}

	e = downloadPictures(0, 3, 0, 3, pathl+fnl, d, date)
	if e != nil {
		return e
	}

	e = mergePicturesAndRotate(pathl, fnl, pathr, fnr,0, 3, 0, 3, 0, color.Black)
	if e != nil {
		return e
	}

	if (setWallpaper) {
		e = setDesktopPicture(pathr + fnr + ".png")
		if e != nil {
			return e
		}
	}
	return nil
}

func once() {
	date, err := getNewestDate()
	if (err != nil) {
		fmt.Println(err.Error())
		return
	}
	if lst_da == date {
		fmt.Printf("Date Same Sleep and Wait for next round\n")
		return
	}

	err = downloadOnce(date, true)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lst_da = date
}

func DownloadYesterday() {
	const TimeFormat = "2006-01-02 15:04:05"

	date, err := getNewestDate()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t, err := time.Parse(TimeFormat, date)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for d := time.Duration(0); d < time.Hour*24; d += time.Minute*10 {
		err := downloadOnce(t.Add(-d).Format(TimeFormat), false)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/ep/tmp", os.ModePerm)

	if err != nil {
		fmt.Println("Cant't mkdir ep/tmp : reason : " + err.Error())
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "yst" {
		DownloadYesterday()
		return
	}
	for {
		once()
		//fmt.Printf("Finished\n")
		time.Sleep(3*time.Minute)
	}
}