require 'json'
require 'net/http'

GIVEN_NAMES = %w{
  Aadya
  Aarav
  Adam
  Aleksandar
  Alexander
  Ali
  Alysha
  Amar
  Amelia
  Anastasia
  Ane
  Anna
  Anya
  Aron
  Aya
  Ayzere
  Bence
  Benjamin
  Camila
  Charlie
  Chloe
  Daniel
  Davit
  Eden
  Elie
  Ellen
  Emil
  Emilija
  Emily
  Emma
  Ethan
  Eva
  Fatima
  Florencia
  Francesco
  Gabriel
  Gabrielle
  George
  Giorgi
  Hanna
  Helena
  Hiro
  Hugo
  Isabella
  Ivaana
  Jack
  Jackson
  Jakub
  James
  Jayden
  Jing
  John
  Joseph
  Junior
  Krishna
  Kseniya
  Laia
  Lamija
  Lana
  Leo
  Liam
  Louise
  Lucas
  Luis
  Luka
  Lukas
  Madison
  Maia
  Maksim
  Manuel
  Marc
  Maria
  Mariam
  Mariami
  Mariana
  Marija
  Martina
  Mary
  Maryam
  Mehdi
  Mehmet
  Miguel
  Min-jun
  Mohamed
  Naranbaatar
  Nathan
  Nikau
  Nikola
  Noah
  Noam
  Noel
  Nora
  Odval
  Oliver
  Olivia
  Oscar
  Paul
  Peter
  Precious
  Rasmus
  Rimas
  Roberts
  Sakura
  Santiago
  Seo-yeon
  Sevinj
  Shristi
  Shu-fen
  Sofia
  Somchai
  Sou
  Stevenson
  Thomas
  Tiare
  Venla
  Victoria
  Viktoria
  Wei
  William
  Ximena
  Yerasyl
  Youssef
  Yusif
  Zahra
  Zeynep
}
# not valid by API checking currently
#  Agustín
#  Ramón
#  Léa
#  Sebastián
#  淑芬
#  Eliška
#  Margrét
#  Sofía
#  서연
#  Μαρία
#  София
#  Георги
#  Ксенія
#  მარიამი
#  Виктория
#  Александр
#  Анастасия

SURNAMES = %w{
  Ahmed
  Andersson
  Andov
  Beridze
  Bianchi
  Borg
  Cohen
  Devi
  Dimitrov
  Gruber
  Hansen
  Horvat
  Hovhannisyan
  Hoxha
  Ivanov
  Jansen
  Joensen
  Johansson
  Johnson
  Jones
  Kazlauskas
  Kim
  Korhonen
  Lee
  Mammadov
  Martin
  Melnyk
  Murphy
  Nagy
  Nielsen
  Novak
  Nowak
  Papadopoulos
  Peeters
  Popa
  Rossi
  Rusu
  Schmit
  Silva
  Singh
  Smirnov
  Smith
  Suzuki
  Tamm
  Tan
  Tjon
  Tremblay
  Wilson
  Wong
  Zogaj
}
# API says not valid for now
#  Horváth
#  Kovačević
#  Məmmədov
#  Jovanović
#  Hernández
#  Hodžić
#  Bērziņš
#  Chén
#  Satō
#  López
#  Müller
#  Nguyễn
#  Novák
#  Rodríguez
#  Wáng
#  Yılmaz
#  佐藤
#  鈴木
#  Иванов
#  Смирно́в
#  Јовановић
#  Παπαδόπουλος
#  Díaz
#  Fernández
#  García
#  González
#

COLORS = %w{
  Red
  Green
  Yellow
  Blue
  Orange
  Purple
  Brown
  Magenta
  Black
  White
}

MONIKERS = %w{
  Street
  Road
  Avenue
  Way
  Circle
  Court
  Parkway
  Alley
  Drive
  Place
  Lane
}
def random_name
  return GIVEN_NAMES.sample, SURNAMES.sample
end

uri = URI('http://localhost:8080/profile')

begin
  req = Net::HTTP::Post.new(uri.path, 'Content-Type' => 'application/json')
  req.body = {name: random_name.join(' '), address: "420 Green Street"}.to_json
  res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    resp = http.request(req)
    puts resp.body
  end
rescue => e
  puts "failed #{e}"
end
