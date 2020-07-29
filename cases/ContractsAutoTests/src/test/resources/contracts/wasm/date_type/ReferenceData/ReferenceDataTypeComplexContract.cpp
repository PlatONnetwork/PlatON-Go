#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;
using namespace std;
/**
 * @author qudong
 * 合约引用类型结合其他类型构造复杂合约
 *
 * */
//基类
class People {
  public:
    void setName(std::string name) {
        this->name = name;
    };
    std::string getName() {
        return this->name;
    }
    void setAge(uint8_t age){
        this->age = age;
    }
    uint8_t getAge(){
        return this->age;
    }
  private:
    std::string name;
    uint8_t age;
    PLATON_SERIALIZE(People, (name)(age))

};
//学生类
class Student : public People {
     private:
        uint64_t sId;//学生编号
        std::array<std::string,10> arrayCourse;//学生课程
        //std::map<std::string,uint64_t> courseScore;//课程对应成绩
        bool sex;//性别--(男生:true、女生:false)
      public:
         void setStudentId(uint64_t sId) {
             this->sId = sId;
         };
         uint64_t getStudentId() {
             return this->sId;
         }
         void setArrayCourse(std::array<std::string,10> arrayCourse){
             this->arrayCourse = arrayCourse;
         }
         std::array<std::string,10> getArrayCourse(){
             return this->arrayCourse;
         }
        /* void setCourseScore(std::string course,uint64_t score){
             this->courseScore[course] = score;
         }
         std::map<std::string,uint64_t>  getCourseScore(){
             return this->courseScore;
         }*/
         void setSex(bool sex) {
             this->sex = sex;
         };
         bool getSex() {
              return this->sex;
         }
         PLATON_SERIALIZE_DERIVED(Student, People, (sId)(arrayCourse)(sex))

};

CONTRACT ReferenceDataTypeComplexContract:public platon::Contract, public Student{

  private:
     platon::StorageType<"a"_n, Student> storage_students;
     platon::StorageType<"b"_n, std::vector<Student>> storage_vector_students;

  public:
    ACTION void init(){  };
    //设置学生基本信息
    ACTION void set_student_id(const uint64_t &sId) {
         storage_students.self().setStudentId(sId);
     };
    CONST uint64_t get_student_id() {
         return  storage_students.self().getStudentId();
     }
     ACTION void set_sex(const bool &sex) {
         storage_students.self().setSex(sex);
     };
     CONST bool get_sex() {
         return  storage_students.self().getSex();
     }
    ACTION void set_name(const std::string &name) {
         storage_students.self().setName(name);
    };
    CONST std::string get_name() {
        return storage_students.self().getName();
    }
    ACTION void set_age(const uint8_t &age){
        storage_students.self().setAge(age);
    }
    CONST uint8_t get_age(){
        return storage_students.self().getAge();
    }
   //设置课程信息
   ACTION void set_array_course(const std::array<std::string,10> &arrayCourse){
       storage_students.self().setArrayCourse(arrayCourse);
    }
   CONST std::array<std::string,10> get_array_course(){
        return storage_students.self().getArrayCourse();
    }
   /*ACTION void set_course_score(const std::string &course,uint64_t &score){
        storage_students.self().setCourseScore(course,score);
    }
   CONST std::map<std::string,uint64_t>  get_course_score(){
        return  storage_students.self().getCourseScore();
    }*/




    //设置学生基本信息
    ACTION void set_student_info(const uint64_t &sId,std::string &name,uint8_t &age,bool &sex){
       set_student_id(sId);
       set_name(name);
       set_age(age);
       set_sex(sex);
   };

    //设置课程成绩
   /* ACTION void set_all_course_score(){
          uint64_t initScore = 60;
          std::array<std::string,10> courseArr = getArrayCourse();
          for(uint8_t i = 0; i <= courseArr.size(); i++ ){
                std::string courseName =  courseArr[i];
                initScore += 10;
                setCourseScore(courseName,initScore);
          }
    };*/
    //修改指定课程成绩
   /* ACTION void set_score_by_coursename(std::string &name){
           std::map<std::string,uint64_t> courseScore = getCourseScore();
            uint64_t initScore = 60;
           if(!courseScore.empty()){
               auto iter = courseScore.find(name);
               if(iter != courseScore.end()){//查询到元素
                    //根据key删除元素
                    getCourseScore().erase(name);
                    //添加元素
                    getCourseScore().insert(pair<std::string,uint64_t>(name,initScore));
               }else{//未查询到元素
                    getCourseScore().insert(pair<std::string,uint64_t>(name,initScore));
               }
           }
    }*/




   //ACTION void set_student_info(const uint64_t &sId,std::string &name,uint8_t &age,bool &sex){


};

PLATON_DISPATCH(ReferenceDataTypeComplexContract,(init)
(set_student_info)
//(set_score_by_coursename)(set_all_course_score)
(set_student_id)(get_student_id)
(set_sex)(get_sex)
(set_name)(get_name)
(set_age)(get_age)
(set_array_course)(get_array_course)
//(set_course_score)(get_course_score)
);


