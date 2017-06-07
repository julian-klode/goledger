# Library written in Go to create (h)ledger files 

This library consists of two components:

1. Parsers for various input formats:
    - CSV files created by acqbanking-cli listtrans
    - CSV files of the Landesbank Berlin (Amazon.de Visa card) 
    - JSON files of the N26 online banking (you could grab them in the inspector)
2. Types and Functions to render hledger files.

Currently not public is the functionality to convert from the parser transactions to the hledger transactions, but the rules for that are rather personal anyway.

# License
Copyright © 2017 Julian Andres Klode

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
