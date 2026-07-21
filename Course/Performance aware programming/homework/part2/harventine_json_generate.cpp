//get the arguments uniform or circular
// seed
// and number of the points needed
// 
// create a json file with the points
// 
// 
#include <algorithm>
#include <cstdlib>
#include<string.h>
#include <iostream>
#include <fstream>
#include <random>
#include <threads.h>
#include <iomanip>

/* Casey's Code */
typedef double f64;
static f64 Square(f64 A)
{
    f64 Result = (A*A);
    return Result;
}

static f64 RadiansFromDegrees(f64 Degrees)
{
    f64 Result = 0.01745329251994329577 * Degrees;
    return Result;
}

// NOTE(casey): EarthRadius is generally expected to be 6372.8
static f64 ReferenceHaversine(f64 X0, f64 Y0, f64 X1, f64 Y1, f64 EarthRadius)
{
    /* NOTE(casey): This is not meant to be a "good" way to calculate the Haversine distance.
       Instead, it attempts to follow, as closely as possible, the formula used in the real-world
       question on which these homework exercises are loosely based.
    */
    
    f64 lat1 = Y0;
    f64 lat2 = Y1;
    f64 lon1 = X0;
    f64 lon2 = X1;
    
    f64 dLat = RadiansFromDegrees(lat2 - lat1);
    f64 dLon = RadiansFromDegrees(lon2 - lon1);
    lat1 = RadiansFromDegrees(lat1);
    lat2 = RadiansFromDegrees(lat2);
    
    f64 a = Square(sin(dLat/2.0)) + cos(lat1)*cos(lat2)*Square(sin(dLon/2));
    f64 c = 2.0*asin(sqrt(a));
    
    f64 Result = EarthRadius * c;
    
    return Result;
}
/* End Casey's Code */

double RandomRange(std::mt19937* gen, double min, double max) {
    std::uniform_real_distribution<double> dist(min, max);
    return dist(*gen);
}

double min(double a, double b) {
    return (a < b) ? a : b;
}

double max(double a, double b) {
    return (a > b) ? a : b;
}

double uniform_generate(int seed, int numberOfPoints){

    std::mt19937 gen(seed);

    long double maxX=180;
    long double maxY=90;
    
    double sum=0;
    double mul = 1/(double)numberOfPoints;
    

    std::ofstream outFile("points.json");
    std::ofstream outFile2("harvestineans.txt");
    
    outFile << std::fixed << std::setprecision(16);
    outFile2 << std::fixed << std::setprecision(16);
    
    outFile << "{" << std::endl;
    outFile << "\"pairs\": [" << std::endl;
    
    for(int i = 0; i < numberOfPoints; i++){
        double x0 = RandomRange(&gen, -maxX, maxX);
        double x1 = RandomRange(&gen, -maxX, maxX);
        double y0 = RandomRange(&gen, -maxY, maxY);
        double y1 = RandomRange(&gen, -maxY, maxY);

        
        double earthRadius = 6372.8;
        double harvestineDistance = ReferenceHaversine(x0, y0, x1, y1, earthRadius);
        sum+=(mul*harvestineDistance);
        outFile2 << harvestineDistance << std::endl;
        
        if(i==numberOfPoints-1)
            outFile << "{" <<"\"x0\": " << x0 << ",\"y0\": " << y0<< ",\"x1\": " << x1<< ",\"y1\": " << y1 << "}" << std::endl;
        else
            outFile << "{" <<"\"x0\": " << x0 << ",\"y0\": " << y0<< ",\"x1\": " << x1<< ",\"y1\": " << y1 << "}," << std::endl;
    }
    
    outFile << "]" << std::endl;
    outFile << "}" << std::endl;
    
    outFile.close();
    outFile2.close();

    return sum;
}



double cluster_generate(int seed, int numberOfPoints){
    std::mt19937 gen(seed);

    int numberofClusters = 1 + (numberOfPoints / 64);
    long double centerx = 0;
    long double centery = 0;
    long double radiusx = 180;
    long double radiusy = 90;
    long double maxX=180;
    long double maxY=90;
    double sum=0;
    double mul = 1/(double)numberOfPoints;
    
    int currentCluster = numberofClusters;
    std::ofstream outFile("points.json");
    std::ofstream outFile2("harvestineans.txt");

    outFile << std::fixed << std::setprecision(16);
    outFile2 << std::fixed << std::setprecision(16);
    
    outFile << "{" << std::endl;
    outFile << "\"pairs\": [" << std::endl;
    
    for(int i = 0; i < numberOfPoints; i++){

        if(currentCluster--==0){
            centerx = RandomRange(&gen, -maxX, maxX);
            centery = RandomRange(&gen, -maxY, maxY);
            radiusx = RandomRange(&gen, 0, maxX);
            radiusy = RandomRange(&gen, 0, maxY);
            currentCluster = numberofClusters;
        }

        double x0 = RandomRange(&gen, max(centerx-radiusx, -maxX), min(centerx+radiusx, maxX));
        double y0 = RandomRange(&gen, max(centery-radiusy, -maxY), min(centery+radiusy, maxY));

        double x1 = RandomRange(&gen, max(centerx-radiusx, -maxX), min(centerx+radiusx, maxX));
        double y1 = RandomRange(&gen, max(centery-radiusy, -maxY), min(centery+radiusy, maxY));
        
        double earthRadius = 6372.8;
        double harvestineDistance = ReferenceHaversine(x0, y0, x1, y1, earthRadius);
        sum+=(mul*harvestineDistance);
        outFile2 << harvestineDistance << std::endl;
        
        if (i==numberOfPoints-1) {
            outFile << "{" <<"\"x0\": " << x0 << ",\"y0\": " << y0<< ",\"x1\": " << x1<< ",\"y1\": " << y1 << "}" << std::endl;
        } else {
            outFile << "{" <<"\"x0\": " << x0 << ",\"y0\": " << y0<< ",\"x1\": " << x1<< ",\"y1\": " << y1 << "}," << std::endl;
        }
    }
    
    outFile << "]" << std::endl;
    outFile << "}" << std::endl;
    
    outFile.close();
    outFile2.close();
    return sum;
}

int main(int argc, char* argv[])
{
    if (argc != 4)
    {
        std::cerr << "Usage: " << argv[0] << " <uniform/cluster> <random seed> <number of pairs to generate>" << std::endl;
        return 1;
    }

    std::string method = argv[1];
    long seed = std::stol(argv[2]);
    long numberOfPoints = std::stol(argv[3]);

    double sum = 0;
    if(method=="uniform"){
        sum = uniform_generate(seed, numberOfPoints);
    }
    else if(method=="cluster"){
        sum = cluster_generate(seed, numberOfPoints);
    }


    std::cout << "Method: " << method << ", Seed: " << seed << ", Number of Points: " << numberOfPoints << ", Sum: " << sum << std::endl;
    
}